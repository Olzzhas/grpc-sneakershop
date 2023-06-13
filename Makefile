sneaker-proto:
	protoc -I .	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	service/sneaker_service/model/v1/sneaker.proto \
	service/sneaker_service/service/v1/sneaker.proto

sneaker-protog:
	protoc -I . --grpc-gateway_out=. \
        --grpc-gateway_opt=logtostderr=true \
        --grpc-gateway_opt=paths=source_relative \
        --grpc-gateway_opt=generate_unbound_methods=true \
        service/sneaker_service/service/v1/sneaker.proto

user-proto:
	protoc -I .	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	service/user_service/model/v1/user.proto \
	service/user_service/service/v1/user.proto

user-protog:
	protoc -I . --grpc-gateway_out=. \
        --grpc-gateway_opt=logtostderr=true \
        --grpc-gateway_opt=paths=source_relative \
        --grpc-gateway_opt=generate_unbound_methods=true \
        service/user_service/service/v1/user.proto


auth-proto:
	protoc -I .	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	service/authentication_service/model/authentication.proto \
	service/authentication_service/service/authentication.proto

auth-protog:
	protoc -I . --grpc-gateway_out=. \
        --grpc-gateway_opt=logtostderr=true \
        --grpc-gateway_opt=paths=source_relative \
        --grpc-gateway_opt=generate_unbound_methods=true \
        service/authentication_service/service/authentication.proto