# Use the Go 1.24 Alpine image as the base image
FROM golang:1.24-alpine

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the application
RUN go build -o main .

# Expose the port that the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
