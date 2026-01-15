# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /sentinel-api cmd/api/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /sentinel-api .
COPY --from=builder /app/web ./web
EXPOSE 3000
CMD ["./sentinel-api"]