# backend

작업 규약은 아래 파일을 따릅니다. 팀원 일부가 Codex를 쓰기 때문에 **같은 파일을 공유합니다.**
규칙을 고칠 일이 있으면 이 파일이 아니라 `AGENTS.md`를 고치고, `make sync-conventions`로
app 리포에도 반영하세요.

@AGENTS.md

## Claude Code 전용 메모

- 커밋 전에 브랜치를 확인합니다. `main`이면 먼저 `git switch -c <type>/<이슈번호>-<요약>`.
- 커밋 메시지 type은 "동작이 바뀌었나"로 고릅니다. 간단한 버그 수정도 `fix`이지 `chore`가 아닙니다.
- `docs/api.md`를 고치지 않고 핸들러의 요청·응답 형태를 바꾸지 않습니다.
- **public 리포입니다.** 시크릿이 들어간 파일을 스테이징하지 마세요.
- 훅이 안 걸리는 것 같으면 `git config core.hooksPath` 값을 확인하세요. `.githooks`여야 합니다.
