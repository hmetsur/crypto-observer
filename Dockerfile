# ---------- build ----------
FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o app ./cmd/app

# ---------- runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/app /app/app
COPY --from=builder /app/configs /app/configs
EXPOSE 8080
# путь к конфигу можно переопределить в docker-compose через ENV CONFIG_PATH
ENV CONFIG_PATH=/app/configs/config.yaml
ENTRYPOINT ["/app/app"]