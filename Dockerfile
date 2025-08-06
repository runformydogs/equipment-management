FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o equipment-manager ./cmd/server

FROM alpine:latest

RUN apk --no-cache add tzdata

COPY --from=builder /app/equipment-manager /equipment-manager

COPY frontend /frontend

EXPOSE 8080

CMD ["/equipment-manager"]