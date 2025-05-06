FROM golang:1.24.2-alpine3.19 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./cmd/backend ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/cmd/backend /app/backend
COPY --from=builder /app/cmd/config.env /app/config.env

EXPOSE 80
CMD ["/app/backend"]