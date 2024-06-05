FROM golang:latest as build
WORKDIR /source
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o mtv_svr ./app

FROM alpine:latest as runtime
WORKDIR /app

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /app/logs
RUN mkdir -p /app/config

COPY --from=build /source/mtv_svr /app
COPY ./config/env.online_docker.toml /app/config/env.toml

ENV APP_ENVIRONMENT development
ENTRYPOINT ["/app/mtv_svr"]
