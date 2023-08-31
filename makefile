swag:
	swag init --parseDependency

up:
	docker-compose up --build

gen_proto:
	protoc --proto_path=proto proto/user.proto --go_out=internal/api --go-grpc_out=internal/api

run:
	go run ./cmd