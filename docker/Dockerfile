FROM golang:alpine as builder 

WORKDIR /build

ADD go.mod .
ADD go.sum .

COPY . .

RUN go build -o app cmd/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app
COPY --from=builder /build/migrations  /build/migrations

COPY config /build/config

CMD ["./app", "--config", "config/deploy.yml"]