#build stage
FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /build

RUN apk add --no-cache git
RUN git clone https://github.com/pressly/goose

WORKDIR /build/goose

RUN go build \
	-tags='no_mysql no_sqlite3 no_ydb no_libsql no_mssql no_vertica' \
	-o goose \
	./cmd/goose

#final stage
FROM alpine:3.21

WORKDIR /root/

COPY --from=builder /build/goose ./

CMD ["./goose"]
