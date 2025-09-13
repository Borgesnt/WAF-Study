FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copia apenas os arquivos de mod primeiro para caching
COPY go.mod go.sum ./
RUN go mod download

# Copia o resto do código
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o waf-proxy .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/waf-proxy .
CMD ["./waf-proxy"]
