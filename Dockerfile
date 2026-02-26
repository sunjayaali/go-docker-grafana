# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app

# ---------- Stage 2: Runtime ----------
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT ["/app/app"]
