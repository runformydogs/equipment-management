FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o equipment-management ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/equipment-management .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./equipment-management"]