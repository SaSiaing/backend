# 가족 가계부 — 작업 규약

가족 구성원이 각자 지출·수입을 기록하고 우리 집 돈이 어디로 얼마나 나가는지
함께 확인하는 Flutter 앱과 Go API 서버입니다.

**규약의 원본은 backend 리포의 이 파일입니다.** Codex가 자동으로 읽습니다.
app 리포에는 `make sync-conventions`로 같은 파일과 규칙 하네스를 복사합니다.

## 저장소

가족 가계부는 리포 두 개로 나뉩니다. 한 창에서 열려면 부모 디렉터리의
`sasiaing.code-workspace` 를 VS Code로 여세요.

| | |
|---|---|
| [SaSiaing/backend](https://github.com/SaSiaing/backend) | Go · Echo · PostgreSQL · sqlc. API 계약(`docs/api.md`)과 작업 순서(`docs/plan.md`)도 여기 |
| [SaSiaing/app](https://github.com/SaSiaing/app) | Flutter |

**두 리포는 브랜치·커밋 규약이 동일합니다.** 규칙을 고칠 때는 backend에서 고치고
`make sync-conventions` 로 app에 복사한 뒤, 양쪽 다 커밋하세요.

## 브랜치

- **`main`에 직접 커밋·push 하지 않습니다.** 항상 feature 브랜치를 팝니다.
- **dev / stage 브랜치는 만들지 않습니다.** MVP 기간에는 `main` 하나로 갑니다.
- 이름: `<type>/<이슈번호>-<요약>`

```
feat/21-transaction-create
fix/24-past-date-reset
refactor/07-repo-interface
```

- 브랜치 번호는 `[04]` 같은 계획 ID가 아니라 **GitHub가 부여한 실제 이슈 번호**를 씁니다.
- 이슈는 backend 리포에 모읍니다. app PR에서는 `SaSiaing/backend#번호`로 연결합니다.
- 원칙은 이슈 하나에 브랜치·PR 하나입니다. 여러 작업을 묶으면 상위 이슈에 체크리스트로 관리합니다.

## 커밋

```
<type>(<scope>): <설명> (#이슈번호)
```

type은 **난이도가 아니라 의도**로 고릅니다. 질문 하나로 갈립니다:

> **동작이 바뀌었나?**

| type | 뜻 |
|---|---|
| `feat` | 기능·코드·파일 추가 (동작 추가) |
| `fix` | 버그 수정 (동작 교정). **난이도 무관** |
| `docs` | 문서 |
| `refactor` | 구조 변경, 코드 제거 (동작 불변) |
| `chore` | 라이브러리, 설정, 코멘트, 포맷 |

scope는 **선택**입니다. 리포가 나뉘어 있어서 `(server)` `(app)` 같은 건 의미가 없습니다.
쓸 거면 리포 안의 영역으로: `(auth)` `(db)` `(chart)` `(ci)`.

```
feat: 내역 등록 엔드포인트 추가 (#21)
fix(auth): 세션 토큰 만료 처리 (#14)
refactor: repo 인터페이스로 service 의존성 분리 (#07)
chore: golangci-lint 1.62 업데이트
```

이 규약은 `.githooks/commit-msg`가 강제합니다. 클론 직후 `make hooks` (backend는 `make dev`)를
한 번 돌려야 활성화됩니다 — git hook은 클론마다 수동 설정이 필요합니다.

## 이슈·PR

- 새 작업과 버그는 backend의 이슈 폼으로 등록합니다. app에는 별도 이슈를 만들지 않습니다.
- PR 제목은 커밋과 같은 `<type>(<scope>): <설명>` 형식을 씁니다.
- PR 본문에 실제 이슈를 연결합니다. 계획 ID만 적지 않습니다.
- 템플릿의 항목을 지우지 말고 해당 없으면 이유를 적습니다.
- `.github/workflows/pr-metadata.yml`이 PR 제목과 이슈 연결을 검사합니다.

## 공통 규칙

- **금액은 정수(엔).** 문자열로 다루지 않고, 부동소수점을 쓰지 않습니다. DB는 `BIGINT`.
- **시각은 전부 `TIMESTAMPTZ`.** API 표현은 RFC3339 UTC (`2026-08-21T09:00:00Z`).
- **집계는 `occurred_at` 기준.** `created_at`이 아닙니다 — 어제 쓴 걸 오늘 넣는 게 일상입니다.
- 삭제는 soft delete (`deleted_at`). 조회는 항상 `deleted_at IS NULL`.
- 에러 응답은 `{ "message": "...", "code": "..." }` 하나로 통일합니다.
- **API를 바꿀 때는 `backend/docs/api.md`를 먼저 고칩니다.** 코드부터 바꾸지 않습니다.

## 범위 밖 — 구현하지 않습니다

| 안 함 | 그래서 없어야 하는 것 |
|---|---|
| 割り勘 / 정산 | `transactions`에 `settled`·`owed`류 컬럼 없음. 사람별 집계에 **차액·잔액·"정산 필요" 표기 금지**. `payer_id`는 **분류 축**이지 채권이 아닙니다 |
| 카드 자동 연동 | 외부 금융 API 없음. CSV 임포트는 v2 |
| 자산 · 투자 | `balances`·`accounts` 테이블 없음. 모든 집계는 기간 합계(flow)이고 누적 잔고(stock)를 계산하지 않습니다 |
| 자동 분류 · 조언 | `category_id`는 언제나 사용자 입력. 추천·경고 문구 없음 |

## 제품 원칙

- **서로 감시하는 느낌이 들지 않게.** "누가 썼는지"는 파악을 위한 것이지 추궁이 아닙니다.
  금지: 랭킹 표기(1위/2위), "가장 많이 쓴 사람", 빨강=과소비 색상. 사람별 색은 채도를 균등하게.
- **입력 마찰을 줄이는 게 최우선.** 목표는 앱 아이콘 탭부터 저장 완료까지 **10초**.
- **과거 날짜 입력이 일상입니다.** 날짜 변경이 번거로우면 앱을 안 씁니다.

## backend 전용

```
cmd/api/          엔트리포인트
internal/handler/ Echo 핸들러
internal/service/ 도메인 로직
internal/repo/    DB 접근 (sqlc 생성 코드 래핑)
db/migrations/    goose
query/            sqlc 입력
docs/             api.md, plan.md
```

- **이 리포는 public입니다.** `.env`, 서비스 계정 키, OAuth 클라이언트 시크릿을 절대 커밋하지 마세요.
  `.gitignore`가 `.env*`를 막고 있지만 `git add -f`로 뚫립니다.
- 안 쓰는 import는 **컴파일 에러**입니다. `organizeImports`를 켜두세요.
- 도메인 에러(`ErrNotFound` / `ErrForbidden` / `ErrConflict`)를 정의하고 `errors.Is`로 분기합니다.
  HTTP 상태 매핑은 `e.HTTPErrorHandler` 한 곳에서만 합니다.
- 500 에러는 `slog`로 로깅하되 클라이언트에는 상세를 숨깁니다.

## app 전용

```
lib/core/       Dio 클라이언트, 라우팅, 테마, 공용 위젯
lib/features/   화면별 (auth / household / transaction / summary / settings)
```

- 패키지: `dio`, `go_router`, `riverpod`, `freezed`, `json_serializable`,
  `flutter_secure_storage`, `google_sign_in`, `share_plus`, `fl_chart`
- 세션 토큰은 `flutter_secure_storage`에 저장합니다. `SharedPreferences` 아닙니다.
- Dio 인터셉터 두 개: Authorization 헤더 자동 첨부, 에러 응답 → 앱 예외 변환.
  화면마다 에러를 파싱하지 않습니다.
- 모델은 freezed + json_serializable로 생성합니다. 손으로 쓰지 마세요.
  생성 파일(`*.g.dart`, `*.freezed.dart`)은 커밋하지 않습니다 — `make build-runner`로 만듭니다.
- **이 리포는 public입니다.** `google-services.json`, `GoogleService-Info.plist`,
  키스토어를 절대 커밋하지 마세요. `.gitignore`가 막고 있지만 `git add -f`로 뚫립니다.

### 화면 만들 때

- 금액 입력은 **숫자 키패드부터**. 금액 → 카테고리 → 저장이 기본 동선입니다.
- 날짜는 기본 오늘이되 **과거 선택이 두 탭 안에** 끝나야 합니다.
- 사람별 색상은 채도를 균등하게. 한 명만 튀면 랭킹처럼 읽힙니다.
- 빈 상태를 항상 만듭니다. 첫 실행에서 흰 화면이 나오면 안 됩니다.
