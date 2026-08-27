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
- `GET /households/:id/transactions` — 필터: `from` `to` `payer_id` `category_id` `type`
- `PATCH /households/:id/transactions/:tid` — 낙관적 락, `version` 불일치 시 409
- `DELETE /households/:id/transactions/:tid`

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
