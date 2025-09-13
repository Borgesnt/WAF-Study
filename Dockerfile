# Usa uma imagem base pequena
FROM golang:1.20-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia apenas os arquivos de mod primeiro para caching
COPY go.mod ./
RUN go mod download

# Copia o resto do código, incluindo o index.html
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o waf-proxy .

# Imagem final, muito mais leve
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/waf-proxy .
# Copia o arquivo do frontend para o container final
COPY index.html .
CMD ["./waf-proxy"]
