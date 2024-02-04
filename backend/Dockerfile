# Use the official Go image as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

RUN go mod download

# Copy the local package files to the container's workspace
COPY . .

# Build the application
RUN go build -o backend .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./backend"]
