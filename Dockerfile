FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o public-api ./cmd/public
RUN go build -o protected-api ./cmd/protected

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates curl

RUN mkdir -p /app/geoip && \
    curl -L "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-Country.mmdb" \
    -o /app/geoip/GeoLite2-Country.mmdb

COPY --from=builder /app/public-api /app/
COPY --from=builder /app/protected-api /app/

EXPOSE 8000

CMD ["/app/public-api"]
