############################
# STEP 1 build executable binary
# This is an intermediate docker image, will be deleted.
############################

FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

ADD . $GOPATH/src/app/
WORKDIR $GOPATH/src/app/

# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -o /go/bin/echo_server

EXPOSE 1323

CMD ["/go/bin/echo_server"]