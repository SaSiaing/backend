package transaction

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
)

func TestTransactionCRUD(t *testing.T) {
	e := newTestServer()
	created := createTransaction(t, e, "/households/home/transactions", map[string]any{
		"payer_id":    "user-1",
		"category_id": "food",
		"type":        "expense",
		"amount":      1200,
		"occurred_at": "2026-09-01T03:00:00Z",
		"memo":        "점심",
	})

	if created.ID == "" || created.Version != 1 {
		t.Fatalf("created transaction = %+v, want id and version 1", created)
	}
	if created.Amount != 1200 || created.Type != TypeExpense {
		t.Fatalf("created transaction = %+v, want amount 1200 and expense", created)
	}

	listRecorder := request(t, e, http.MethodGet, "/households/home/transactions?type=expense", nil)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d: %s", listRecorder.Code, http.StatusOK, listRecorder.Body.String())
	}
	var listed []Transaction
	decodeJSON(t, listRecorder, &listed)
	if len(listed) != 1 || listed[0].ID != created.ID {
		t.Fatalf("listed transactions = %+v, want created transaction", listed)
	}

	updatePath := "/households/home/transactions/" + created.ID
	updateRecorder := request(t, e, http.MethodPatch, updatePath, map[string]any{
		"amount":  1500,
		"memo":    "점심 수정",
		"version": 1,
	})
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d: %s", updateRecorder.Code, http.StatusOK, updateRecorder.Body.String())
	}
	var updated Transaction
	decodeJSON(t, updateRecorder, &updated)
	if updated.Amount != 1500 || updated.Memo != "점심 수정" || updated.Version != 2 {
		t.Fatalf("updated transaction = %+v, want updated fields and version 2", updated)
	}

	conflictRecorder := request(t, e, http.MethodPatch, updatePath, map[string]any{
		"amount":  2000,
		"version": 1,
	})
	if conflictRecorder.Code != http.StatusConflict {
		t.Fatalf("conflict status = %d, want %d", conflictRecorder.Code, http.StatusConflict)
	}

	getRecorder := request(t, e, http.MethodGet, updatePath, nil)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getRecorder.Code, http.StatusOK)
	}

	deleteRecorder := request(t, e, http.MethodDelete, updatePath, nil)
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteRecorder.Code, http.StatusNoContent)
	}

	deletedRecorder := request(t, e, http.MethodGet, updatePath, nil)
	if deletedRecorder.Code != http.StatusNotFound {
		t.Fatalf("deleted get status = %d, want %d", deletedRecorder.Code, http.StatusNotFound)
	}

	emptyListRecorder := request(t, e, http.MethodGet, "/households/home/transactions", nil)
	if got := strings.TrimSpace(emptyListRecorder.Body.String()); got != "[]" {
		t.Fatalf("list after delete = %s, want []", got)
	}
}

func TestCreateTransactionRejectsInvalidAmount(t *testing.T) {
	e := newTestServer()
	recorder := request(t, e, http.MethodPost, "/households/home/transactions", map[string]any{
		"payer_id":    "user-1",
		"category_id": "food",
		"type":        "expense",
		"amount":      0,
		"occurred_at": time.Now().UTC(),
	})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	var response errorResponse
	decodeJSON(t, recorder, &response)
	if response.Code != "invalid_argument" {
		t.Fatalf("error code = %q, want invalid_argument", response.Code)
	}
}

func newTestServer() *echo.Echo {
	e := echo.New()
	handler := NewHandler(NewService(NewMemoryRepository()))
	handler.Register(e)
	return e
}

func createTransaction(t *testing.T, e *echo.Echo, path string, body any) Transaction {
	t.Helper()
	recorder := request(t, e, http.MethodPost, path, body)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	var created Transaction
	decodeJSON(t, recorder, &created)
	return created
}

func request(t *testing.T, e *echo.Echo, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *strings.Reader
	if body == nil {
		requestBody = strings.NewReader("")
	} else {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		requestBody = strings.NewReader(string(encoded))
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	recorder := httptest.NewRecorder()
	e.ServeHTTP(recorder, req)
	return recorder
}

func decodeJSON(t *testing.T, recorder *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.Unmarshal(recorder.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
}
