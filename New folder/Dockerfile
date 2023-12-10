# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy all the files from the current directory to the /app directory inside the container
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go application and name the executable as "app"
RUN go build -o app

# Expose port 8000 for the application to listen on
EXPOSE 8000

# Command to run the executable when the container starts
CMD ["./app"]
