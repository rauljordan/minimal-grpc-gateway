go:
	go build -o main .
protos:
	protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:. \
	  proto/api/v1/api.proto
	protoc -I. --go_out=plugins=grpc:${GOPATH}/src ./proto/api/v1/api.proto
