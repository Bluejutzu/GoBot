FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o gobot .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gobot .
COPY --from=builder /app/.env .

RUN chmod +x /app/gobot

CMD ["/app/gobot"]
