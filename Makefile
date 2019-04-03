SHELL := /bin/bash
PROTO_FILES := $(wildcard protobuf/*.proto)
PROTO_DEFS  := $(PROTO_FILES:.proto=.pb.go)

.PHONY: proto2
proto2: $(PROTO_FILES)

protobuf/%.pb.go: protobuf/%.proto
	protoc --go_out=plugins=grpc:. $<

.PHONY: clean
clean:
	rm -f protobuf/*.pb.go
