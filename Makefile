PROTO_FILES := $(ls protobuf/*.proto)
PROTO_DEFS  := $(PROTO_FILES:.proto=.pb.go)

.PHONY: proto
proto: $(PROTO_DEFS)

protobuf/%.pb.go: protobuf/%.proto
	protoc --go_out=plugins=grpc:. $<

.PHONY: clean
clean:
	rm -f protobuf/*.pb.go
