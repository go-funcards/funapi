# FunCards RestAPI

![License](https://img.shields.io/dub/l/vibe-d.svg)

## Install Plugins

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/swaggo/swag/cmd/swag@v1.8.3
```

## Generate RSA Keys

```shell
mkdir -p .data/cert
openssl genrsa -out .data/cert/id_rsa 4096
openssl rsa -in .data/cert/id_rsa -pubout -out .data/cert/id_rsa.pub
```

## Run:

```shell
docker compose up -d
```

## License

Distributed under MIT License, please see license file within the code for more details.
