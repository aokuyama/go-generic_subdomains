FROM golang:1.17.2-alpine3.14

ARG app_dir="/app"

RUN mkdir -p $app_dir
WORKDIR $app_dir

COPY ./api ${app_dir}/api
WORKDIR ${app_dir}/api
RUN go mod download

COPY ./param ${app_dir}/param
WORKDIR ${app_dir}/param
RUN go mod download

WORKDIR $app_dir
