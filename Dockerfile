# Base the image off of the lastest version of 1.18
FROM golang:1.18

# Create container directory structure
RUN mkdir -p /home/app

# Copy files from the host to the container (in the specified directory)
COPY . /home/app

# Set working directory
WORKDIR /home/app

# Build the application
RUN go build -v -o ./main ./cmd/services

# Run the app in the image
CMD ["./main"]
