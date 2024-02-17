# create binary from official golang image
FROM golang:1.21-bookworm as builder

# create and change to app dir
WORKDIR /app

# copy over mod and sum file
COPY go.* ./
RUN go mod download && go mod verify

# copy all files in local dir into container
COPY . ./

# build the binary
RUN go build -v -o gin-server ./main/main.go

# debian slim image for lean prod container
FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# copy binary to prod image from builder stage
COPY --from=builder /app/gin-server /app/gin-server

# EXPOSE 80


# Run the web service on container startup
CMD ["/app/gin-server"]

