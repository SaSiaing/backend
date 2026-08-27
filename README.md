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

[AGENTS.md](AGENTS.md) 에 있습니다. Codex는 이 파일을 자동으로 읽고, Claude Code는 `CLAUDE.md`가 임포트합니다.

- `main`에 직접 커밋·push 하지 않습니다. `<type>/<이슈번호>-<요약>` 브랜치를 팝니다
- dev / stage 브랜치는 만들지 않습니다 (MVP 기간)
- 커밋: `<type>(<scope>): <설명> (#이슈)` — type은 **"동작이 바뀌었나"**로 고릅니다
  `feat` `fix` `docs` `refactor` `chore` · scope는 선택

`.githooks/`가 강제합니다. 사람이 치든 Claude가 치든 Codex가 치든 똑같이 걸립니다.
**규칙을 고쳤으면 `make sync-conventions` 로 app 리포에도 반영하세요.**
