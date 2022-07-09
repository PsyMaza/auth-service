GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

SERVICE_NAME=auth_service
SERVICE_PATH=g6834/team17/api

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

.PHONY: run
run:
	go run cmd/app/main.go

.PHONY: lint
lint:
	golangci-lint run --timeout 5m --config .golangci.yaml -v ./...

.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'

# ----------------------------------------------------------------

.PHONY: generate
generate: .generate-install-buf .generate-go .generate-finalize-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go:
	$(BUF_EXE) generate

.generate-finalize-go:
	mv pkg/$(SERVICE_NAME)/gitlab.com/$(SERVICE_PATH)/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	rm -rf pkg/$(SERVICE_NAME)/gitlab.com/
	cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init gitlab.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

# ----------------------------------------------------------------

.PHONY: deps
deps: deps-go

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

# ----------------------------------------------------------------

.PHONY: build
build: generate .build


.build:
	go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-v -o ./bin/auth-service$(shell go env GOEXE) ./cmd/app/main.go

# ----------------------------------------------------------------

docs:
	docker build --tag swaggo/swag:1.8.1 . --file swaggo.Dockerfile && \
	docker run --rm --volume ${PWD}:/app --workdir /app swaggo/swag:1.8.1 /root/swag init \
		--parseDependency \
		--parseInternal \
		--dir ./internal/api \
		--generalInfo swagger.go \
		--output ./api/swagger/public \
		--parseDepth 1
