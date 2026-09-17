# API 계약

> **이 문서를 먼저 고치고 코드를 고칩니다.** 순서가 반대면 프론트·백이 갈라집니다.
> 변경 시 PR 체크리스트에 "api.md 갱신"을 넣으세요.

## 공통 규칙

| 항목 | 규칙 |
|---|---|
| 금액 | **정수** (엔 단위). 문자열·소수 아님. `1200` |
| 날짜·시각 | RFC3339 UTC. `2026-08-21T09:00:00Z` |
| 인증 | `Authorization: Bearer <token>` |
| 에러 | `{ "message": "...", "code": "..." }` |
| 페이지네이션 | 커서 방식. `?cursor=<opaque>&limit=50` (offset 쓰지 않음) |
| 삭제 | soft delete. 조회 결과에 나오지 않음 |

## 엔드포인트

### 상태
- `GET /health` — 인증 없이 PostgreSQL 연결 상태 확인. 정상 시 `{ "status": "ok" }`

### 인증
- `POST /auth/google` — body `{ "id_token": "..." }` → 세션 토큰

### 가계
- `POST /households`
- `GET /households`
- `PATCH /households/:id` — `month_start_day` 변경 (#34)
- `POST /households/:id/invitations`
- `POST /invitations/:code/accept`

### 카테고리
- `GET /households/:id/categories`
- `POST /households/:id/categories`
- `PATCH /households/:id/categories/:cid`
- `DELETE /households/:id/categories/:cid`

### 내역
- `POST /households/:id/transactions`
  - 요청: `{ "payer_id": "user-1", "category_id": "food", "type": "expense", "amount": 1200, "occurred_at": "2026-09-01T03:00:00Z", "memo": "점심" }`
  - 응답: 생성된 내역과 `version: 1`
- `GET /households/:id/transactions` — 필터: `from` `to` `payer_id` `category_id` `type`
- `GET /households/:id/transactions/:tid`
- `PATCH /households/:id/transactions/:tid`
  - 요청: 바꿀 필드와 현재 `version` (`{ "amount": 1500, "version": 1 }`)
  - 응답: 수정된 내역과 증가한 `version`; 버전 불일치 시 409
- `DELETE /households/:id/transactions/:tid`

> 현재 첫 CRUD 단계에서는 내역을 메모리에 저장합니다. 서버를 재시작하면 초기화되며 인증·멤버십 검증, 커서 페이지네이션, DB migration/sqlc는 후속 작업입니다.

### 집계
- `GET /households/:id/summary/by-category?from=&to=`
- `GET /households/:id/summary/monthly?months=6`
- `GET /households/:id/summary/by-payer?from=&to=`

---

## TODO — #02에서 채울 것

- [ ] 각 엔드포인트의 요청·응답 JSON 예시
- [ ] 에러 `code` 목록 확정 (`not_found` / `forbidden` / `conflict` / `invalid_argument`)
- [ ] 팀 전원 리뷰 후 머지

> 완벽할 필요 없습니다. 바뀌는 걸 전제로 하되, **바뀔 때 문서를 먼저 고치는 습관**이 목적입니다.
