FROM golang:alpine

RUN apk add build-base
RUN apk add --no-cache git
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/myapp
ADD . /go/src/myapp
WORKDIR /go/src/myapp

RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

EXPOSE 80

ENTRYPOINT ["/go/bin/myapp"]