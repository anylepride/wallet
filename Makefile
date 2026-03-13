.PHONY: init

init:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: generate

generate:
	protoc -I ./proto \
      --go_out ./proto --go_opt paths=source_relative \
      --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
      --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
      --openapiv2_out ./proto --openapiv2_opt logtostderr=true \
      proto/wallet/*.proto