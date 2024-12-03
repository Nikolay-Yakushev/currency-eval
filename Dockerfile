FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main /app/cmd/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/internal/repository/pg/migrations /app/internal/repository/pg/migrations

ENTRYPOINT ["./main"]