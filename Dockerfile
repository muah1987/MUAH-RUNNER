# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o muah-runner ./cmd/muah-runner/

# Stage 2: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/muah-runner .
EXPOSE 8080 9090
ENTRYPOINT ["./muah-runner"]
CMD ["start"]
