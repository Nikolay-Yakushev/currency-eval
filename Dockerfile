FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main /currency_app/cmd/currency_app/main.go /currency_app/cmd/currency_app/config.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/internal/repository/postgres/migrations /app/internal/repository/postgres/migrations

ENTRYPOINT ["./main"]