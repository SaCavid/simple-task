FROM golang:1.15.8-alpine

RUN apk add build-base
RUN apk add --no-cache git

ADD . /go/src/myapp
WORKDIR /go/src/myapp

RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

EXPOSE 80

ENTRYPOINT ["/go/bin/myapp"]