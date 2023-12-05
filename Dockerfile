FROM golang:1.16

WORKDIR /go/src/github.com/ietxaniz/delock

# Copy your project files into the container
COPY . .

