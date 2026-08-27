#!/bin/sh
# main 브랜치에서 커밋을 시도하면 경고. (실제 차단은 .githooks/pre-commit)
branch=$(git symbolic-ref --short HEAD 2>/dev/null) || exit 0
[ "$branch" = "main" ] && echo "  main입니다. 먼저 'git switch -c <type>/<번호>-<요약>' 하세요." >&2
exit 0
