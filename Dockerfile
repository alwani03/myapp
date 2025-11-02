FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/myapp /app/myapp
EXPOSE 8080
CMD ["/app/myapp"]