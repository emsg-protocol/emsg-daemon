# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o emsg-daemon ./cmd/daemon/main.go

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/emsg-daemon ./emsg-daemon
#COPY --from=builder /app/docs ./docs
COPY --from=builder /app/internal ./internal
COPY --from=builder /app/api ./api
#COPY --from=builder /app/config ./config
COPY --from=builder /app/go.mod ./go.mod
COPY --from=builder /app/go.sum ./go.sum
EXPOSE 8080
CMD ["./emsg-daemon"]
