package main

import app "github.com/go-funcards/funapi/cmd"

//go:generate curl -LJO --output-dir proto/user_service/v1 https://raw.githubusercontent.com/go-funcards/user-service/main/proto/v1/user.proto
//go:generate curl -LJO --output-dir proto/board_service/v1 https://raw.githubusercontent.com/go-funcards/board-service/main/proto/v1/board.proto
//go:generate curl -LJO --output-dir proto/tag_service/v1 https://raw.githubusercontent.com/go-funcards/tag-service/main/proto/v1/tag.proto
//go:generate curl -LJO --output-dir proto/category_service/v1 https://raw.githubusercontent.com/go-funcards/category-service/main/proto/v1/category.proto
//go:generate curl -LJO --output-dir proto/card_service/v1 https://raw.githubusercontent.com/go-funcards/card-service/main/proto/v1/card.proto
//go:generate curl -LJO --output-dir proto/authz_service/v1 https://raw.githubusercontent.com/go-funcards/authz-service/main/proto/v1/checker.proto
//go:generate curl -LJO --output-dir proto/authz_service/v1 https://raw.githubusercontent.com/go-funcards/authz-service/main/proto/v1/subject.proto

//go:generate protoc -I proto --go_out=./proto/user_service/v1 --go-grpc_out=./proto/user_service/v1 proto/user_service/v1/user.proto
//go:generate protoc -I proto --go_out=./proto/board_service/v1 --go-grpc_out=./proto/board_service/v1 proto/board_service/v1/board.proto
//go:generate protoc -I proto --go_out=./proto/tag_service/v1 --go-grpc_out=./proto/tag_service/v1 proto/tag_service/v1/tag.proto
//go:generate protoc -I proto --go_out=./proto/category_service/v1 --go-grpc_out=./proto/category_service/v1 proto/category_service/v1/category.proto
//go:generate protoc -I proto --go_out=./proto/card_service/v1 --go-grpc_out=./proto/card_service/v1 proto/card_service/v1/card.proto
//go:generate protoc -I proto --go_out=./proto/authz_service/v1 --go-grpc_out=./proto/authz_service/v1 proto/authz_service/v1/checker.proto proto/authz_service/v1/subject.proto

//go:generate swag init --parseDependency -g cmd/serve.go

func main() {
	app.Execute()
}
