FROM golang:1.21-bullseye AS builder

WORKDIR /app

# SQLite ke liye gcc chahiye
RUN apt-get update && apt-get install -y gcc libsqlite3-dev && rm -rf /var/lib/apt/lists/*

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o hinamd .

# ── Final image ──
FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/hinamd .

EXPOSE 8080

CMD ["./hinamd"]
