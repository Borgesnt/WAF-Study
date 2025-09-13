FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o waf-proxy .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/waf-proxy .
CMD ["./waf-proxy"]
