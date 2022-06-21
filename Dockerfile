# Base the image off of the lastest version of 1.18
FROM golang:1.18

# Make the detsination directory
RUN mkdir /app

# Copy files to the new directory
ADD . /app

# Set working directory
WORKDIR /app

# Build the application
RUN go build -o main ./cmd/services

# Run the app in the image
CMD ["/app/main"]