MIGRATE_DB=postgres://postgres:postgres@localhost:5435/postgres?sslmode=disable

migrate-up:
	goose -dir migrations postgres "$(MIGRATE_DB)" up

migrate-down:
	goose -dir migrations postgres "$(MIGRATE_DB)" down

run:
	go run cmd/auth/main.go

lint:
	golangci-lint run --config=.golangci.yml ./...

fmt:
	gci write -s standard -s default -s "prefix(auth)" .