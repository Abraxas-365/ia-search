# Start from a base Go image
FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

ENV DOCKERIZE_VERSION v0.6.1
RUN apk add --no-cache openssl \
    && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
# Set the working directory
RUN mkdir /app
WORKDIR /app

# Copy the entire project into the working directory
COPY . .

# Download the dependencies
# RUN go mod download
RUN go mod tidy

# Build the application
RUN go build cmd/main.go

# Expose the port the application will run on
EXPOSE 8080

# Run the application
# CMD ["./main"]
CMD dockerize -wait tcp://rabbitmq:15672 -timeout 60m ./main
