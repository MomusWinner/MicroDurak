#build stage
FROM golang:1.24.1-alpine3.21 AS builder
ARG SERVICE_NAME

WORKDIR /micro-durak/$SERVICE_NAME

ADD ./cmd ./
ADD ./services ./
ADD ./lib ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o app cmd/$SERVICE_NAME/main.go

#final stage
FROM alpine:3.21

WORKDIR /root/

COPY --from=builder /micro-durak/$SERVICE_NAME/app ./

CMD ["./app"]

EXPOSE 8080
