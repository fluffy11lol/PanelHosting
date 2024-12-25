FROM golang:1.23-alpine AS builder

WORKDIR /service/app


COPY /services/authentication-microservice .

COPY /pkg/logger ../../pkg/logger


RUN ls -l

RUN go mod download

WORKDIR /service/app/cmd/main

RUN go build -o ./authentication-service


FROM alpine:latest

WORKDIR /app/configs

COPY --from=builder /service/app/configs .

WORKDIR /app

COPY --from=builder /service/app/cmd/main/authentication-service .

EXPOSE 8081
EXPOSE 9091

CMD ["./authentication-service"]