FROM golang:alpine3.18 AS builder

COPY . /Fitness_REST_API/

WORKDIR /Fitness_REST_API/

RUN go mod download

RUN go test ./internal/handler
RUN go test ./internal/repository/

RUN go build -o ./bin/app cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 Fitness_REST_API/bin/app .
COPY --from=0 Fitness_REST_API/configs configs/

EXPOSE 80

CMD ["./app"]