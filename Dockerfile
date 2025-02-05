# Use official Golang image as the base image
FROM golang:1.23.4

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container (including the services folder)
COPY . .

# Build the Go app (main.go is inside services, so we build from there)
RUN GOOS=linux go build -o /app/binarygo ./services/cmd/main.go

# Expose the port the app will run on
EXPOSE 8080

# Define the command to run the app
ENTRYPOINT [ "/app/binarygo" ]
