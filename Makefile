run:
	go run cmd/auth/main.go

lint:
	golangci-lint run --config=.golangci.yml ./...

fmt:
	gci write -s standard -s default -s "prefix(auth)" .