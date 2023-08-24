FROM golang:1.21-alpine3.18

WORKDIR /app

COPY . .

# Download dependencies
RUN go mod download

RUN go build -o create-tables ./scripts

CMD ["./create-tables"]
