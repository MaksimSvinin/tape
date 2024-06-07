project-name=tape
build_filename=./build/$(project-name)
go_ldflags="-w -s"

# Выполнить проверку линтером
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.1
# doc: https://golangci-lint.run/usage/install/
.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: generate
generate:
	go generate ./...
