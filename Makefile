GOLANGCI_LINT := v1.39.0

.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint $(GOLANGCI_LINT)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANGCI_LINT)

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

.PHONY: run
run:
	@go run main/main.go