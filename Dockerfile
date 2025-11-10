FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/auth ./cmd/auth


FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/auth .
COPY .env .

EXPOSE 8080
ENTRYPOINT ["./auth"]
