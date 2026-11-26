FROM golang:1.19 AS builder
RUN mkdir /datahub-api
ADD . /datahub-api
WORKDIR /datahub-api
ARG PATUSER
ARG PATTOKEN

# Setting up PAT for private repos
RUN git config --global url."https://${PATUSER}:${PATTOKEN}@dev.azure.com".insteadOf "https://dev.azure.com"

# Go env
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go env -w CGO_ENABLED="0"
RUN go env -w GOPRIVATE=dev.azure.com/conductortech
RUN go env -w GO111MODULE='on'

# Fetching dependencies
RUN go mod download

# Getting CLI tools
RUN go install github.com/swaggo/swag/cmd/swag@v1.6.5
RUN go install github.com/gobuffalo/packr/packr@v1.30.1
RUN swag -v
RUN packr version
# RUN version -v

# Build
RUN swag init --output cmd/http/docs --generalInfo cmd/http/main.go
RUN packr -v
RUN go build -o bin/datahub-api -ldflags="-s -w" ./cmd/http


FROM alpine:latest
RUN mkdir -p /go/bin
WORKDIR /go/bin
COPY --from=builder /datahub-api/bin/datahub-api .
EXPOSE 8080
RUN mkdir /logs
RUN chmod -R 755 /logs
ENTRYPOINT /go/bin/datahub-api