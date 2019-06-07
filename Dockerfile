# Start from golang v1.11 base image
FROM golang:1.11 as builder

# Add Maintainer Info
LABEL maintainer="werty1st <werty1st@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/werty1st/sparkgi

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY main.go .

# Download dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/sparkgi .



######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/sparkgi .

EXPOSE 1025

CMD ["./sparkgi"] 