PROTO_SRC = api/commands.proto
PROTO_OUT = .

.PHONY: all proto clean test bench

all: proto

proto:
	protoc --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
	       --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
	       $(PROTO_SRC)

clean:
	rm -f $(PROTO_OUT)/api/*.pb.go

test:
	go test ./server -v

bench:
	go test -bench . ./benchmarks
