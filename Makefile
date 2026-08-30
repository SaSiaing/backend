.PHONY: hooks dev down test lint sqlc migrate seed sync-conventions

## hooks: git hook 활성화 (클론 직후 1회, dev가 자동으로 부름)
hooks:
	@git config core.hooksPath .githooks
	@echo "git hooks 활성화됨 → .githooks"

## dev: 로컬 개발환경 기동
dev: hooks
	docker compose up -d
	@echo "postgres  localhost:5432  (kakeibo / kakeibo / kakeibo_dev)"

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

## sync-conventions: AGENTS.md와 Git/Codex 하네스를 app 리포로 복사
sync-conventions:
	@./scripts/sync-conventions.sh
