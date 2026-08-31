# 백엔드 폴더 구조

Echo는 폴더 구조를 강제하지 않습니다. 이 프로젝트는 기능이 커질 때 관련 코드를 찾기 쉽도록 **기능별 패키지**를 기본으로 사용합니다. PostgreSQL 연결이나 sqlc 생성 코드처럼 여러 기능이 공유하는 인프라만 별도로 둡니다.

```text
backend/
├── cmd/
│   └── api/
│       └── main.go              # 실행 진입점과 의존성 조립
├── internal/
│   ├── database/
│   │   └── postgres.go          # PostgreSQL 연결 풀
│   ├── health/
│   │   ├── handler.go           # Health 라우트
│   │   └── handler_test.go
│   ├── <feature>/               # auth, household, transaction 등
│   │   ├── handler.go           # HTTP 요청·응답과 라우트
│   │   ├── service.go           # 비즈니스 규칙이 있을 때 추가
│   │   ├── repository.go        # DB 접근
│   │   └── model.go             # 해당 기능의 타입
│   └── db/                      # sqlc 생성 코드, 직접 수정 금지
├── db/
│   └── migrations/              # goose migration
├── query/                       # sqlc 입력 SQL
├── docs/
├── docker-compose.yml
└── go.mod
```

## 요청 흐름

```text
클라이언트 → handler → service → repository → sqlc → PostgreSQL
```

- 단순 CRUD처럼 비즈니스 규칙이 없으면 service를 생략할 수 있습니다.
- handler는 SQL을 직접 실행하지 않습니다.
- repository는 Echo의 Context나 HTTP 상태 코드를 사용하지 않습니다.
- sqlc 생성 파일은 직접 고치지 않고 `query/`와 `sqlc.yaml`을 수정한 뒤 다시 생성합니다.
- `cmd/api/main.go`는 객체를 조립하고 라우트를 등록하지만 CRUD 로직은 갖지 않습니다.

## 전역 폴더를 먼저 만들지 않는 이유

- 최상위 `models/` 대신 `internal/<feature>/model.go`에 기능 타입을 둡니다.
- `pkg/`는 다른 프로젝트에서도 가져다 쓸 공개 라이브러리가 실제로 생길 때만 추가합니다.
- 모든 라우트를 모은 `router/` 대신 각 기능의 handler가 자신의 라우트를 등록합니다.

새 기능은 `internal/<feature>/`에서 시작하고, 둘 이상의 기능이 실제로 공유하게 된 코드만 공통 패키지로 옮깁니다.
