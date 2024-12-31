# Use Golang as the base image
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy the app's source code
COPY . .

# Download dependencies and build the app
RUN go mod tidy
RUN go build -o app .

# Expose the port the app runs on
EXPOSE 8080

# Start the app
CMD ["./app"]
