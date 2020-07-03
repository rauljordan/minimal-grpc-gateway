# Minimal gRPC Gateway Example

This PR creates a sample gRPC-JSON gateway using Go.

### Running:

```build
git clone https://github.com/rauljordan/minimal-grpc-gateway && cd minimal-grpc-gateway
go build . -o main && ./main
```

This will spin up a grpc server on port 4000 and a grpc gateway on 8080.

```text
INFO[0000] gRPC server listening on port                 address="127.0.0.1:8080"
INFO[0000] Starting gRPC gateway                         address="127.0.0.1:4000" prefix=gateway
```

### Interacting with the gateway

```build
curl -X POST http://localhost:8080/api/v1/users/signup -H "accept: application/json" -d '{"username": "someone", "password": "1234}'
```

Expected Behavior:
```build
{"jwtKey":"aGVsbG8td29ybGQ="}
```

### Regenerating protobufs

Install protoc [here](https://google.github.io/proto-lens/installing-protoc.html), then

```text
make
```