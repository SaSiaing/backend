#!/usr/bin/env bash
# 라벨 8개 · 마일스톤 4개 · 이슈 39개 생성.
#
#   ./scripts/gh-bootstrap.sh [owner/repo] <W1시작일 YYYY-MM-DD>
#   예) ./scripts/gh-bootstrap.sh 2026-08-31
#
# 이슈·마일스톤은 backend 리포에 모읍니다. GitHub 마일스톤은 리포 단위라
# 두 리포에 나누면 4개를 양쪽에 똑같이 만들어야 합니다. area:fe 이슈도 여기 둡니다.
#
# 이슈 제목의 [21] [23a] 접두사는 docs/plan.md의 ID입니다.
# GitHub이 붙이는 번호와 다를 수 있으니 문서 참조는 접두사를 기준으로 하세요.
set -euo pipefail

if [ $# -eq 1 ]; then REPO="SaSiaing/backend"; START="$1"
else REPO="${1:?owner/repo 또는 시작일을 넘기세요}"; START="${2:?W1 시작일(YYYY-MM-DD)을 넘기세요}"; fi

d() { date -j -v+"$1"d -f "%Y-%m-%d" "$START" +%Y-%m-%dT23:59:59Z; }

echo "==> 라벨"
lbl() { gh label create "$1" -R "$REPO" -c "$2" -d "$3" --force >/dev/null && echo "    $1"; }
lbl "area:be"     "1D76DB" "Go 백엔드"
lbl "area:fe"     "0E8A16" "Flutter"
lbl "area:infra"  "5319E7" "배포·CI·개발환경"
lbl "area:docs"   "BFD4F2" "API 계약·문서"
lbl "blocker"     "B60205" "이게 안 끝나면 다른 사람이 멈춤"
lbl "good-first"  "7057FF" "Go/Flutter 처음이어도 잡을 만한 것"
lbl "must-have"   "D93F0B" "제품 우선순위 직결 — 버릴 후보에서 제외"
lbl "v2"          "FEF2C0" "v1 완료 후 백로그"

echo "==> 마일스톤"
ms() { gh api "repos/$REPO/milestones" -f title="$1" -f due_on="$2" -f description="$3" >/dev/null 2>&1 \
       && echo "    $1 (~$(echo "$2" | cut -dT -f1))" || echo "    $1 (이미 있음)"; }
ms "M1 뼈대" "$(d 7)" \
"리포·개발환경·배포 파이프라인·Go/Flutter 뼈대. OAuth 클라이언트 발급까지."
ms "M2 입력" "$(d 14)" \
"인증, 가계 생성, 내역 등록·목록·삭제, 입력 화면. 여기서부터 실제로 지출을 기록하기 시작."
ms "M3 집계·공유" "$(d 21)" \
"카테고리별→월별→사람별 집계, 대시보드, 초대 코드, 내역 수정."
ms "M4 정착" "$(d 28)" \
"폴링 동기화, 로딩·에러·빈 상태 정리, 실사용에서 나온 마찰 제거."

echo "==> 이슈"
iss() { # milestone labels title body
  gh issue create -R "$REPO" -m "$1" -l "$2" -t "$3" -b "$4" >/dev/null && echo "    $3"
}

# ---------------- M1 ----------------
M="M1 뼈대"
iss "$M" "area:docs,blocker" "[02] API 계약 초안" \
'**모든 작업의 출발점.** 이게 나와야 Flutter 담당이 목 데이터로 선행할 수 있습니다.

- [ ] 엔드포인트 목록 (인증 / 가계 / 내역 / 집계)
- [ ] 요청·응답 JSON 예시
- [ ] 공통 규칙: 금액 정수(엔), RFC3339 UTC, 에러 `{message, code}`, `Authorization: Bearer`
- [ ] 에러 code 목록 확정
- [ ] 팀 전원 리뷰 후 머지

완벽할 필요 없습니다. 바뀔 때 **문서를 먼저 고치는 습관**이 목적입니다.'
iss "$M" "area:infra,blocker" "[01] 리포지토리 · 브랜치 전략" \
'- [x] 모노레포 `/server` `/app` `/docs`
- [ ] GitHub ruleset — main 직접 push 차단, PR 1인 리뷰
- [x] `.gitignore`, `.githooks/`, `AGENTS.md`
- [x] README 로컬 실행 절차

로컬 훅은 `--no-verify`로 뚫립니다. **진짜 방어선은 GitHub ruleset입니다.**'
iss "$M" "area:infra,good-first" "[03] 개발환경 통일" \
'- [x] `docker-compose.yml` — Postgres 16
- [x] `Makefile` — dev / test / sqlc / migrate / seed / hooks
- [ ] `.vscode/settings.json` (formatOnSave, organizeImports, golangci-lint)
- [ ] `.vscode/extensions.json` — `golang.go`, `Dart-Code.flutter`
- [ ] `.golangci.yml` — errcheck, govet, staticcheck, ineffassign
- [ ] 3명 모두 `make dev`로 DB 뜨는 것 확인

`organizeImports` 자동화가 중요합니다. Go는 안 쓰는 import가 **컴파일 에러**입니다.'
iss "$M" "area:be,blocker" "[04] Go 프로젝트 초기화 + Echo 세팅" \
'- [ ] `go mod init`
- [ ] `cmd/api`, `internal/handler`, `internal/service`, `internal/repo`
- [ ] Echo + Logger·Recover 미들웨어
- [ ] CORS (Flutter 웹 디버그용)
- [ ] `GET /health` → `{"status":"ok"}`
- [ ] 환경변수 로딩 (`DATABASE_URL`, `PORT`)

학습: `http.Handler`, 미들웨어 체인'
iss "$M" "area:be,blocker" "[05] 공통 에러 핸들러 · 응답 포맷" \
'- [ ] `e.HTTPErrorHandler` → `{message, code}` 통일
- [ ] 도메인 에러 타입 (`ErrNotFound`, `ErrForbidden`, `ErrConflict`)
- [ ] 도메인 에러 → HTTP 상태 매핑
- [ ] 500은 `slog`로 로깅, 클라이언트에는 상세 숨김

W1에 해두면 Flutter 에러 처리가 한 곳에서 끝납니다. 나중에 하면 핸들러 전부 수정입니다.
학습: `errors.Is` / `errors.As`'
iss "$M" "area:be,blocker" "[06] 스키마 DDL + 마이그레이션" \
'```
households          id, name, month_start_day, created_at
household_members   household_id, user_id, display_name, role, joined_at
users               id, google_sub, email, name, created_at
categories          id, household_id, name, kind, sort_order, deleted_at
transactions        id, household_id, payer_id, category_id, type, amount,
                    occurred_at, memo, version, created_at, updated_at, deleted_at
invitations         code, household_id, expires_at, used_at
```

- [ ] goose 도입 + 첫 마이그레이션
- [ ] `amount`는 `BIGINT` (엔 단위 정수)
- [ ] 시각은 전부 `TIMESTAMPTZ`
- [ ] ID는 UUID
- [ ] soft delete `deleted_at`
- [ ] 인덱스 `(household_id, occurred_at)`, `(household_id, payer_id)`

`month_start_day`를 미리 둡니다. 나중에 넣으면 집계 쿼리를 전부 다시 씁니다.'
iss "$M" "area:be" "[07] sqlc 세팅 + 첫 쿼리" \
'- [ ] `sqlc.yaml`
- [ ] `query/`에 첫 쿼리
- [ ] `make sqlc` 생성 확인
- [ ] repo 인터페이스 정의

학습: 코드 생성 워크플로, Go 인터페이스'
iss "$M" "area:infra,blocker" "[08] 배포 파이프라인" \
'**기능이 하나도 없을 때 배포부터 뚫습니다.** 마지막에 하면 반드시 늦어집니다.

- [ ] Dockerfile (멀티스테이지, distroless)
- [ ] Cloud Run 배포
- [ ] Cloud SQL 또는 Supabase Postgres 연결
- [ ] `min-instances=1`, `timeout=3600s`, CPU always allocated
- [ ] GitHub Actions: main push → 자동 배포
- [ ] 배포 URL로 `/health` 성공

이틀 이상 막히면 **Fly.io로 전환.** Cloud Run 설정에 매몰되지 않기.'
iss "$M" "area:infra,blocker" "[12] GCP OAuth 클라이언트 발급" \
'**Phase 2에서 W1으로 당겼습니다.** 코드 의존성이 없고 3명 SHA-1을 모아야 해서,
W2 첫날 시작하면 그 주 인증 작업이 통째로 밀립니다. **반나절 잡으세요.**

- [ ] GCP 프로젝트 + OAuth 동의 화면
- [ ] Android 클라이언트 ID (SHA-1 debug/release 각각)
- [ ] iOS 클라이언트 ID (번들 ID)
- [ ] Web 클라이언트 ID (서버 검증용 audience)
- [ ] 3명의 개발 머신 SHA-1 전부 등록'
iss "$M" "area:fe,blocker" "[09] Flutter 프로젝트 초기화" \
'- [ ] `flutter create`
- [ ] `dio`, `go_router`, `riverpod`, `freezed`, `json_serializable`, `flutter_secure_storage`
- [ ] `lib/features/`, `lib/core/`
- [ ] 3명 모두 에뮬레이터 실행 확인'
iss "$M" "area:fe" "[10] Flutter API 클라이언트 뼈대" \
'- [ ] Dio 인스턴스 + baseUrl 환경 분리 (로컬/배포)
- [ ] 인터셉터: Authorization 헤더 자동 첨부
- [ ] 인터셉터: 에러 응답 → 앱 예외 변환
- [ ] `docs/api.md` 기준 freezed 모델
- [ ] `/health` 호출로 연결 확인

[02]만 끝나면 백엔드 완성 전에도 목 데이터로 진행 가능합니다.'
iss "$M" "area:fe,good-first" "[11] Flutter 화면 라우팅 뼈대" \
'- [ ] go_router 설정
- [ ] 빈 화면: 로그인 / 대시보드 / 내역목록 / 내역입력 / 설정
- [ ] 하단 네비게이션
- [ ] 로그인 여부 리다이렉트 골격'

# ---------------- M2 ----------------
M="M2 입력"
iss "$M" "area:be" "[13] Go: 구글 ID 토큰 검증 + 유저 upsert" \
'- [ ] `POST /auth/google` — body `{ "id_token": "..." }`
- [ ] `google.golang.org/api/idtoken`으로 검증
- [ ] `payload.Subject`(sub)로 유저 식별 — **이메일로 하지 않기**
- [ ] users upsert
- [ ] 자체 세션 토큰 발급 후 반환'
iss "$M" "area:be,blocker" "[14] Go: 세션 토큰 발급 · 인증 미들웨어" \
'- [ ] 자체 JWT 발급 (유효기간 7일)
- [ ] 미들웨어: `Authorization: Bearer` 파싱 → 검증 → context에 userID
- [ ] context 헬퍼 (`SetUserID`/`GetUserID`) — 타입 안전하게
- [ ] 미인증 401

구글 ID 토큰은 1시간이면 만료되니 그대로 세션으로 쓰지 않습니다.
학습: `context.Context`, 미들웨어 합성'
iss "$M" "area:fe" "[15] Flutter: 구글 로그인" \
'- [ ] `google_sign_in` 연동
- [ ] ID 토큰 → `/auth/google`
- [ ] 세션 토큰을 `flutter_secure_storage`에 저장
- [ ] 앱 시작 시 자동 로그인
- [ ] 로그아웃'
iss "$M" "area:be" "[16] Go: 가계 생성 · 조회" \
'- [ ] `POST /households` — 생성자가 owner
- [ ] `GET /households` — 내가 속한 가계 목록
- [ ] 생성 시 기본 카테고리 자동 삽입 (식비/일용품/교통 등)'
iss "$M" "area:be,blocker" "[17] Go: 멤버십 권한 미들웨어" \
'**나중에 넣으면 모든 핸들러를 다시 만집니다. 반드시 여기서.**

- [ ] `/households/:id/*`에서 멤버십 확인
- [ ] 비멤버 403
- [ ] context에 householdID + role 저장
- [ ] role 기반 제한 (owner만 멤버 삭제 등)'
iss "$M" "area:be,good-first" "[20] Go: 카테고리 CRUD" \
'- [ ] 목록 / 추가 / 이름변경 / 삭제(soft)
- [ ] 수입용·지출용 구분 (`kind`)'
iss "$M" "area:be,must-have" "[21] Go: 내역 등록" \
'- [ ] `POST /households/:id/transactions`
- [ ] 검증: amount > 0, type 유효, category가 같은 가계 소속인지
- [ ] `payer_id` 기본값 = 요청자 (다른 멤버 대신 입력도 허용)
- [ ] 201 + 생성된 리소스 반환'
iss "$M" "area:be" "[22] Go: 내역 목록 조회" \
'- [ ] `GET /households/:id/transactions`
- [ ] 필터: 기간, payer, category, type
- [ ] 정렬 `occurred_at DESC`
- [ ] **커서** 페이지네이션 (offset 말고)
- [ ] `deleted_at IS NULL` 필수'
iss "$M" "area:be" "[23a] Go: 내역 삭제 (soft delete)" \
'[23]에서 쪼갠 앞부분. **W2 말부터 실사용을 시작하므로 오입력을 되돌릴 수단이 그날부터 필요합니다.**
혼자 쓰는 동안은 지우고 다시 넣으면 충분합니다.

- [ ] `DELETE /households/:id/transactions/:tid` — soft delete
- [ ] 권한: 멤버 전원 (한 주머니 전제 + 감시 느낌 회피)

낙관적 락(PATCH)은 [23b]로 분리 — 둘이 같은 걸 만질 때 물건입니다.'
iss "$M" "area:be,good-first" "[37] make seed — 집계 검증용 더미 내역" \
'W3에 월별 추이(6개월)를 만들 때 실데이터는 며칠치뿐입니다. **시드 없이는 [29]를 검증할 방법이 없습니다.**

- [ ] `cmd/seed` — 가계 1개, 멤버 3명, 6개월치 내역
- [ ] 카테고리·지불자·금액 분포를 그럴듯하게
- [ ] `make seed`'
iss "$M" "area:fe,must-have" "[24] Flutter: 내역 입력 화면" \
'**실제로 가장 많이 쓰는 화면입니다. 입력 마찰을 줄이는 데 시간을 쓰세요.**

- [ ] 금액 키패드 (숫자 입력 최적화, 금액부터)
- [ ] 수입/지출 토글
- [ ] 카테고리 선택
- [ ] 날짜 (기본 오늘, **과거 입력 필수**)
- [ ] 지불자 선택 (기본 본인)
- [ ] 메모

목표는 앱 아이콘 탭부터 저장 완료까지 **10초**. 화면 진입 후가 아니라 아이콘부터입니다.'
iss "$M" "area:fe" "[25] Flutter: 내역 목록 화면" \
'- [ ] 날짜별 그룹핑
- [ ] 무한 스크롤
- [ ] 사람별 색상·아바타 구분 (채도 균등 — [36] 참고)
- [ ] 당겨서 새로고침
- [ ] 빈 상태 UI'
iss "$M" "area:fe" "[26a] Flutter: 스와이프 삭제" \
'[26]에서 쪼갠 앞부분.

- [ ] 스와이프 삭제 + 확인 다이얼로그
- [ ] 삭제 후 목록 갱신

수정 UI와 409 충돌 처리는 [26b].'

# ---------------- M3 ----------------
M="M3 집계·공유"
iss "$M" "area:be,must-have" "[28] Go: 카테고리별 집계" \
'**M3에서 가장 먼저 만듭니다.** 혼자 쓴 며칠치로도 값이 나오는 유일한 집계입니다.

- [ ] `GET /households/:id/summary/by-category?from=&to=`
- [ ] 금액 내림차순
- [ ] 비율(%) 포함
- [ ] `occurred_at` 기준 (created_at 아님)'
iss "$M" "area:be,must-have" "[29] Go: 월별 수지 추이" \
'- [ ] `GET /households/:id/summary/monthly?months=6`
- [ ] 월별 수입·지출·수지
- [ ] `month_start_day` 반영
- [ ] 데이터 없는 달도 0으로 채우기

대시보드 상단 "이번 달 수입·지출·수지"는 별도 엔드포인트 없이 **이 응답의 현재 달 행**을 씁니다.'
iss "$M" "area:be,must-have" "[27] Go: 사람별 집계" \
'**주 후반에 합니다.** 초대([18][19])가 끝나야 값이 의미를 가집니다 — 혼자면 100% 하나만 찍힙니다.

- [ ] `GET /households/:id/summary/by-payer?from=&to=`
- [ ] 지출·수입 각각 합계
- [ ] 타임존 `AT TIME ZONE ''Asia/Tokyo''`
- [ ] `occurred_at` 기준

**차액·잔액·"정산 필요" 표기 금지.** `payer_id`는 분류 축이지 채권이 아닙니다.
학습: SQL 집계, 타임존'
iss "$M" "area:be" "[18] Go: 초대 코드" \
'- [ ] `POST /households/:id/invitations` — 코드 생성 (만료 24h)
- [ ] `POST /invitations/:code/accept` — 참여
- [ ] 만료·사용완료 코드 거부
- [ ] 이미 멤버면 409'
iss "$M" "area:fe" "[19] Flutter: 가계 생성 · 참여 화면" \
'- [ ] 가계 없을 때: 생성 / 코드로 참여 선택
- [ ] 초대 코드 표시 + 공유 (share_plus)
- [ ] 가계 전환 (여러 개 속할 경우)

**M2가 아니라 M3인 이유:** 텅 빈 앱에 가족을 초대하면 한 번 열고 안 옵니다.
보여줄 게 생긴 다음에 부릅니다.'
iss "$M" "area:fe,must-have" "[30] Flutter: 대시보드" \
'**여기가 앱의 존재 이유입니다.** 세 화면이 서로 다른 질문에 답합니다.

- [ ] 이번 달 수입·지출·수지 요약 ([29] 응답의 현재 달 행)
- [ ] 사람별 지출 — "누가 무엇을 썼는지"
- [ ] 카테고리별 (도넛) — "어디에 나가는지"
- [ ] 월별 추이 — "늘고 있는지 줄고 있는지"
- [ ] 월 전환
- [ ] `fl_chart` 도입'
iss "$M" "area:be" "[23b] Go: 내역 수정 (낙관적 락)" \
'- [ ] `PATCH /households/:id/transactions/:tid` — 부분 업데이트
- [ ] `version` 불일치 시 409
- [ ] 권한: **멤버 전원** (대신 입력을 허용하면서 수정만 막는 건 앞뒤가 안 맞습니다)

학습: 동시성 제어, 부분 업데이트'
iss "$M" "area:fe" "[26b] Flutter: 내역 수정 + 충돌 처리" \
'- [ ] 항목 탭 → 수정 화면
- [ ] 409 충돌 시 안내 후 재조회'
iss "$M" "area:be,area:fe" "[34] 월 시작일 설정" \
'**컬럼은 [06]에 있는데 사용자가 바꿀 화면이 어디에도 없습니다.**
원칙 5 — 25일 급여 기준으로 한 달을 세는 집이 있습니다.

- [ ] `PATCH /households/:id` — `month_start_day`
- [ ] 설정 화면에 입력 UI
- [ ] [29] 집계가 값을 반영하는지 확인'
iss "$M" "area:be,area:fe" "[35] 카테고리 사용 빈도순 정렬" \
'기능 B "카테고리는 **자주 쓰는 순으로**". [20]의 `sort_order`는 수동 정렬이라 다른 물건입니다.
**10초 목표에 직접 기여합니다.**

- [ ] 최근 N건 기준 사용 빈도 집계
- [ ] 입력 화면 카테고리 목록을 빈도순으로'
iss "$M" "area:docs,area:fe" "[36] 중립 카피 · 색상 가이드" \
'원칙 2 — **서로 감시하는 느낌이 들지 않게.** 지금 [25][30]에는 "사람별 색상 구분"만 있어서,
구현하다 보면 자연히 랭킹처럼 됩니다. **금지 목록이 필요합니다.**

- [ ] 금지: 랭킹 표기(1위/2위), "가장 많이 쓴 사람", 빨강=과소비
- [ ] 금지: 차액·잔액·"정산 필요" — `payer_id`는 분류 축이지 채권이 아닙니다
- [ ] 사람별 색상은 채도 균등 (한 명만 튀지 않게)
- [ ] [25] [30] 리뷰 체크리스트로 사용'

# ---------------- M4 ----------------
M="M4 정착"
iss "$M" "area:fe" "[31] Flutter: 폴링 동기화" \
'- [ ] 앱 포그라운드 복귀 시 재조회
- [ ] 화면 진입 시 재조회
- [ ] 필요 시 주기 폴링 (30초)

가족은 동시 편집을 거의 안 해서 이걸로 체감상 충분합니다. SSE는 v2.'
iss "$M" "area:fe" "[32] 로딩 · 에러 · 빈 상태 정리" \
'- [ ] 전 화면 로딩 인디케이터
- [ ] 네트워크 에러 → 재시도 버튼
- [ ] 401 → 로그인 화면으로
- [ ] 빈 상태 안내 문구'
iss "$M" "" "[33] 2주 실사용에서 나온 마찰 제거" \
'**"실사용 테스트"가 아닙니다.** 실사용은 W2 말에 이미 시작했습니다. M4는 테스트가 아니라 수리입니다.

W4 초에 열어서 그 주에 채워 넣는 상위 이슈입니다.

- [ ] 2주간 쓰면서 걸린 것들을 개별 이슈로 등록
- [ ] 그중 W4 안에 고칠 수 있는 것 처리'

echo
echo "완료. https://github.com/$REPO/milestones"
