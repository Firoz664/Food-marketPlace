# Use the official Golang image to create a build artifact.
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Compile the Go app for Linux (in case you're building on a non-linux system)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest  

# Add CA Certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /app/main .

# Expose port 9001 to the outside world
EXPOSE 9000

# Command to run the executable
CMD ["./main"]
