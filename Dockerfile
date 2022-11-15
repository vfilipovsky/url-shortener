FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o .bin/server cmd/server/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/.bin/server .
COPY --from=builder /app/.env .
COPY --from=builder /app/public /app/public

CMD ["./server"]