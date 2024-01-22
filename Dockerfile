# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder

RUN apk add git

RUN addgroup -g 1000 -S appgroup && \
  adduser -u 1000 -S appuser -G appgroup

RUN mkdir app

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY . ./

# Download all the dependencies
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 go build -o /app/main .

RUN chown -R appuser:appgroup /app

USER 1000

EXPOSE 3000

ENTRYPOINT [ "/app/main" ]

