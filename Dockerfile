FROM golang:1.21-alpine AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS webhook-receiver
CMD ["go", "run", "./cmd/webhook_receiver/main.go"]

FROM base AS hinatazaka-blog-notifier
RUN apk add --no-cache bash
COPY ./cmd/hinatazaka_blog_notifier/script.sh /script.sh
RUN chmod +x /script.sh
CMD ["/script.sh"]

FROM base AS create-tables
CMD ["go", "run", "./scripts/create_tables.go"]
