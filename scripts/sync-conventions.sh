#!/usr/bin/env sh
# 커밋 규약 훅을 app 리포로 복사합니다.
# 두 리포가 같은 부모 디렉터리에 있어야 합니다 (워크스페이스 기본 배치).
#
# 규칙을 고쳤으면 backend 에서 이걸 돌리고, app 쪽도 커밋하세요.
set -e
DEST="${1:-../app}"

if [ ! -d "$DEST/.githooks" ]; then
  echo "  $DEST/.githooks 가 없습니다. app 리포 경로를 인자로 넘기세요:" >&2
  echo "    ./scripts/sync-conventions.sh ../app" >&2
  exit 1
fi

cp .githooks/commit-msg .githooks/pre-commit .githooks/pre-push "$DEST/.githooks/"
chmod +x "$DEST/.githooks/"*
echo "  .githooks → $DEST/.githooks 복사 완료"
echo "  app 리포에서 커밋하는 걸 잊지 마세요."
