# Official Golang image
FROM golang:1.23-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Arguments
ARG PORT

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ./*.go ./
COPY ./docs ./docs
COPY ./database/*.go ./database/
COPY ./handlers/*.go ./handlers/
COPY ./models/*.go ./models/
COPY ./workers/*.go ./workers/

# Build the Go app
RUN go build -o main .

# Expose port to the outside world
EXPOSE ${PORT}

# Command to run the executable
CMD ./main
