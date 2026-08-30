# backend

가족 가계부 API 서버. Go · Echo · PostgreSQL · sqlc

프런트엔드는 [SaSiaing/app](https://github.com/SaSiaing/app) 입니다.
두 리포를 한 창에서 열려면 부모 디렉터리의 `sasiaing.code-workspace` 를 여세요.

## 로컬 실행

```
cp .env.example .env    # 값 채우기
make dev                # Postgres 기동 + git hook 활성화
go run ./cmd/api
```

`make dev`를 **반드시 한 번은 돌려야 합니다.** git hook은 클론마다 수동 활성화가 필요하고,
`make dev`가 그걸 대신해 줍니다. (`make hooks`만 따로 돌려도 됩니다)

## 문서

- [docs/api.md](docs/api.md) — **API 계약. 코드보다 먼저 고칩니다.**
- [docs/plan.md](docs/plan.md) — v1 작업 순서

## 작업 규약

[AGENTS.md](AGENTS.md)에 있습니다. Codex가 프로젝트를 열 때 자동으로 읽습니다.

- `main`에 직접 커밋·push 하지 않습니다. `<type>/<이슈번호>-<요약>` 브랜치를 팝니다
- dev / stage 브랜치는 만들지 않습니다 (MVP 기간)
- 커밋: `<type>(<scope>): <설명> (#이슈)` — type은 **"동작이 바뀌었나"**로 고릅니다
  `feat` `fix` `docs` `refactor` `chore` · scope는 선택

`.githooks/`가 강제하므로 사람이 직접 작업하든 Codex를 쓰든 똑같이 적용됩니다.
**규칙을 고쳤으면 `make sync-conventions` 로 app 리포에도 반영하세요.**

## Codex 팀 설정

- `AGENTS.md` — 제품 원칙, 코드 규칙, 브랜치·커밋 규약. Codex가 자동으로 읽습니다.
- `.codex/config.toml` — workspace-write 샌드박스, 요청 기반 승인, 세션 시작 브랜치 확인.
- `.codex/rules/` — GitHub·Docker의 읽기 전용 명령 허용 목록.
- `.githooks/` — `main` 직접 커밋·push와 잘못된 커밋 메시지를 실제로 차단합니다.
- `.github/` — 이슈·PR 템플릿과 PR 제목·이슈 연결 검사.

처음 열 때 프로젝트를 **trusted**로 승인하고 `/hooks`에서 프로젝트 훅을 검토·승인하세요.
Codex의 작업 루트는 부모 디렉터리가 아니라 `backend` 또는 `app` 저장소로 선택해야 합니다.
`make dev`를 한 번 실행하면 Git 훅까지 활성화됩니다.
