FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/...

# ---- runtime ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server ./server
# config/ must live at /app/config/ so runtime.Caller path resolution works correctly.
# Sensitive values (JWT keys, STEAM_API_KEY) are injected via environment variables at runtime.
COPY --from=builder /app/config/ ./config/

EXPOSE 8080

ENTRYPOINT ["./server", "-env", "prod"]
