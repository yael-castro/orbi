.PHONY: grpc consumer proto

grpc:
	sh scripts/build.bash grpc

consumer:
	sh scripts/build.bash consumer

proto:
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/messages.proto proto/service.proto