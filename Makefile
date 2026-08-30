.PHONY: hooks dev run down test lint sqlc migrate seed sync-conventions

## hooks: git hook 활성화 (클론 직후 1회, dev가 자동으로 부름)
hooks:
	@git config core.hooksPath .githooks
	@echo "git hooks 활성화됨 → .githooks"

## dev: 로컬 개발환경 기동
dev: hooks
	docker compose up -d
	@echo "postgres  localhost:5432  (kakeibo / kakeibo / kakeibo_dev)"

## run: 로컬 DB와 API 서버 실행 (.env 필요)
run: dev
	@test -f .env || { echo ".env가 없습니다. cp .env.example .env를 먼저 실행하세요."; exit 1; }
	@set -a; . ./.env; set +a; go run ./cmd/api

down:
	docker compose down

test:
	go test ./...

lint:
	golangci-lint run

sqlc:
	sqlc generate

migrate:
	goose -dir db/migrations postgres "$$DATABASE_URL" up

seed:
	go run ./cmd/seed

## sync-conventions: AGENTS.md와 Git/Codex/PR 하네스를 app 리포로 복사
sync-conventions:
	@./scripts/sync-conventions.sh
