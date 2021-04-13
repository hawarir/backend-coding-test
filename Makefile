GOLANGCI_LINT := v1.39.0

.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint $(GOLANGCI_LINT)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANGCI_LINT)

.PHONY: loadtest-prepare
loadtest-prepare:
	@echo "Installing hey"
	@GO111MODULE=off go get github.com/rakyll/hey

.PHONY: lint
lint:
	@golangci-lint run \
		--enable=gochecknoglobals \
		--enable=gochecknoinits \
		--enable=goconst \
		--enable=gocyclo \
		--enable=golint \
		--enable=unconvert \
		./...

.PHONY: test
test:
	@go test -v -coverprofile=coverage.out ./...

.PHONY: loadtest
loadtest:
	@hey -n 10000 -c 50 -m GET http://localhost:8010/rides	

.PHONY: run
run:
	@DB_PATH=./rides.db PORT=8010 go run main/main.go
