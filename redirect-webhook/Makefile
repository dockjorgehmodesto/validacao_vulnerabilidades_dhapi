install:
	go env
	git config --global url."git@ssh.dev.azure.com:v3/conductortech/billing/middleware-golang".insteadOf https://dev.azure.com/conductortech/billing/_git/middleware-golang
	go env -w GOPROXY=https://proxy.golang.org
	go env -w CGO_ENABLED="0"
	go env -w GOPRIVATE=dev.azure.com/conductortech
	go env -w GO111MODULE='on'
	go env -w GOBIN=$GOPATH/bin
	go mod download
	go go get github.com/swaggo/swag/cmd/swag@v1.6.5
	go get -u github.com/gobuffalo/packr/packr
	swag -v
	packr version
	version -v

pack:
	packr -v

build:
	swag init init --output cmd/http/docs --generalInfo cmd/http/main.go
	go build -o bin/datahub-api -ldflags="-s -w" ./cmd/http

run:
	go run cmd/http/main.go

all: install pack build
