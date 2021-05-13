FROM golang:1.15.8-alpine

RUN apk add --no-cache git

ADD . /go/src/myapp
WORKDIR /go/src/myapp
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

#RUN go get myapp
#RUN go install
EXPOSE 80
ENTRYPOINT ["/go/bin/myapp"]