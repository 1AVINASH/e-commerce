# Stage 1: Build the Go app
FROM golang:1.23.4 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main .
# Use a debugging step to confirm the file exists
RUN ls -la

# Stage 2: Run the app
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /src/main .

EXPOSE 8080
RUN echo "This message will appear during the Docker build process."
RUN ls -la .
CMD ["./main"]
