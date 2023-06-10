proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	service/sneaker_service/service/sneaker.proto \
	service/sneaker_service/model/sneaker.proto