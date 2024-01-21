# Build stage
FROM golang:1.21-alpine AS builder

ENV USER=builder

RUN apk add --no-cache build-base

RUN adduser -D -g "" $USER 

WORKDIR /app

RUN chown -R root:root /app


COPY --chown=root:$USER go.mod go.sum ./
COPY --chown=root:$USER . .

RUN go mod tidy
RUN go mod download

# Build the main file
RUN GOOS=linux go build -o main cmd/inariam/main.go 

# Run stage
FROM alpine:latest


ENV USER=runner
COPY --from=builder /app /app
# COPY --chown=root:$USER config.yml /app/config.yml

RUN adduser -D -g "" $USER

WORKDIR /app

# Setup stage

CMD ["./main", "server"]
