# Defining App builder image
FROM golang:alpine AS builder

# Add git to determine build git version
RUN apk add --no-cache --update git

# Set GOPATH to build Go app
ENV GOPATH=/go

# Set apps source directory
ENV SRC_DIR=${GOPATH}/src/go-rest-api

# Define current working directory
WORKDIR ${SRC_DIR}

# Copy apps source code to the image
COPY . ${SRC_DIR}

# Build App
RUN ./build.sh

# Defining App image
FROM alpine:latest

RUN apk add --no-cache --update ca-certificates

# Copy App binary to image
COPY --from=builder /go/bin/go-rest-api /usr/local/bin/go-rest-api
COPY config.yml config.yml
COPY prod.env prod.env

EXPOSE 8000

ENTRYPOINT ["go-rest-api"]
CMD ["serve-rest", "-c", "config.yml", "-e", "prod"]