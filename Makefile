include config.env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD)"


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.33.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.33.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7



generate: 
	make generate-wallet-api 
	make generate-openapi
	mkdir -p pkg/swagger
	$(LOCAL_BIN)/statik -src=pkg/swagger -dest=pkg -include="*.css,*.html,*.js,*.json,*.png" -f

generate-wallet-api:
	mkdir -p pkg/wallet_v1
	protoc --proto_path api/wallet_v1 --proto_path vendor.protogen \
	--go_out=pkg/wallet_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/wallet_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/wallet_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	api/wallet_v1/wallet.proto
 

# OpenAPI is generated once, in a single protoc invocation, so allow_merge
# actually merges wallet_v1  one api.swagger.json.
generate-openapi: generate-wallet-api  
	mkdir -p pkg/swagger
	protoc --proto_path api --proto_path vendor.protogen \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/wallet_v1/wallet.proto  

server-up:
	go run cmd/grpc_server/main.go

 
goose-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v
goose-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

goose-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

goose-create:
	mkdir -p $(LOCAL_MIGRATION_DIR)
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create $(NAME) sql

lint:
	$(LOCAL_BIN)/golangci-lint run ./...


 
tests-cover:	
	go test -count=5 \
		-coverpkg=github.com/pavel-pj/wallet/internal/api/wallet,github.com/pavel-pj/wallet/internal/service/wallet \
		-coverprofile=coverage.out \
		./internal/api/wallet/tests ./internal/service/wallet/tests
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html	
 

lib:
	GOPROXY=direct go get github.com/pavel-pj/platform_lib@latest


vendor-proto:
	@if [ ! -f vendor.protogen/buf/validate/validate.proto ]; then \
		rm -rf vendor.protogen/protovalidate &&\
		git clone --depth=1 https://github.com/bufbuild/protovalidate vendor.protogen/protovalidate &&\
		mkdir -p vendor.protogen/buf/validate &&\
		cp vendor.protogen/protovalidate/proto/protovalidate/buf/validate/*.proto vendor.protogen/buf/validate/ &&\
		rm -rf vendor.protogen/protovalidate ;\
	fi	
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -f vendor.protogen/protoc-gen-openapiv2/options/annotations.proto ]; then \
		rm -rf vendor.protogen/openapiv2 &&\
		git clone --depth=1 --branch v2.33.0 https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi

 

DC = docker compose --env-file config.env -f compose.dev.yaml

dev-build:
	$(DC) up --build --no-attach postgres
dev-up:
	$(DC) up --no-attach postgres
dev-down:
	$(DC) down -v
dev-logs:
	$(DC) logs -f app

loadtest:
	go run ./cmd/loadtest   -wallet 1c860e04-d620-4e3d-9b25-c3c0edbdf465   -rate 3000   -duration 20s   -workers 300   -op DEPOSIT   -amount 1

