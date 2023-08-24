# Build stage
FROM golang:1.21-alpine3.18 AS build

WORKDIR /app

COPY . .

# Download dependencies and build the application
RUN go mod download
RUN go build -o create-tables ./scripts

# Runtime stage
FROM alpine:3.18

# If there are any runtime dependencies, install them here

WORKDIR /app

# Copy only the compiled application from the build stage
COPY --from=build /app/create-tables ./create-tables

CMD ["./create-tables"]
