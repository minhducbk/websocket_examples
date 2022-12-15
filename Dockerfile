# Build stage
FROM golang:1.19.3-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN go build -o main ws/main.go

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8084
CMD ["/app/main"]