FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o messenger-backend .

# ───────────────────────────────────────────

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/messenger-backend .
COPY --from=builder /wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 8080

CMD ["./messenger-backend"]
