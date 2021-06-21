.PHONY: build
build: vendor-proto .generate .build

PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-snippet-api
		protoc -I vendor.protogen \
				--go_out=pkg/ocp-snippet-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-snippet-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-snippet-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--validate_out lang=go:pkg/ocp-snippet-api \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/ocp-snippet-api/ocp-snippet-api.proto
		mv pkg/ocp-snippet-api/github.com/ozoncp/ocp-snippet-api/pkg/ocp-snippet-api/* pkg/ocp-snippet-api/
		rm -rf pkg/ocp-snippet-api/github.com
		mkdir -p cmd/ocp-snippet-api

PHONY: .build
.build:
		CGO_ENABLED=0 GOOS=linux go build -o bin/ocp-snippet-api cmd/ocp-snippet-api/main.go

PHONY: vendor-proto
vendor-proto: .vendor-proto

PHONY: .vendor-proto
.vendor-proto:
		mkdir -p vendor.protogen
		mkdir -p vendor.protogen/api/ocp-snippet-api
		cp api/ocp-snippet-api/ocp-snippet-api.proto vendor.protogen/api/ocp-snippet-api
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/github.com/envoyproxy &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate ;\
		fi

lint:
	golangci-lint run