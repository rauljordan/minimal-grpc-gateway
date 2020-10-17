# Minimal gRPC Gateway Example

This PR creates a sample gRPC-JSON gateway using Go.

### Fetching

To use in your projects, you can use go modules and fetch the library with
```
go get github.com/rauljordan/minimal-grpc-gateway
```

### Running the Example

```build
git clone https://github.com/rauljordan/minimal-grpc-gateway && cd minimal-grpc-gateway/example
make go
./main
```

This will spin up a grpc server on port 4000 and a grpc gateway on 8080.

```text
INFO[0000] gRPC server listening on port                 address="localhost:4500"
INFO[0000] Starting gRPC gateway                         address="localhost:5000" prefix=gateway
```

### Interacting with the gateway

```build
curl -X POST http://localhost:5000/api/v1/users/signup -H "accept: application/json" -d '{"username": "someone", "password": "1234}'
```

Expected Behavior:
```build
{"jwtKey":"aGVsbG8td29ybGQ="}
```

### Regenerating protobufs

Install protoc [here](https://google.github.io/proto-lens/installing-protoc.html), then

```text
cd example && make protos
```
