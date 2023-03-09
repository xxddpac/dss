APP_NAME=dss

.PHONY: all build swag
all:
	build
build:
	mkdir -p bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/${APP_NAME}
swag:
	swag init -o core/docs