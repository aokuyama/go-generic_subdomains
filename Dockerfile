FROM golang:1.17.2-alpine3.14
RUN go install github.com/golang/mock/mockgen@v1.6.0

ARG app_dir="/app"

RUN mkdir -p $app_dir
WORKDIR $app_dir

COPY ./go.mod $app_dir
COPY ./go.sum $app_dir

RUN go mod download
