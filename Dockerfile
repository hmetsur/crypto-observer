# ---------- build stage ----------
FROM golang:1.24-alpine AS builder


ENV CGO_ENABLED=0

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -trimpath -ldflags="-s -w" -o app ./cmd


FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /app/app /app/app


EXPOSE 8080


ENV DB_DSN="postgres://user:pass@db:5432/crypto?sslmode=disable"


ENTRYPOINT ["/app/app"]