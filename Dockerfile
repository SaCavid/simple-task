FROM golang:alpine
RUN apk add --no-cache git
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 80
CMD ["/app/main"]


#FROM golang:1.15.8-alpine
#
#RUN apk add --no-cache git
#
#ENV CGO_ENABLED=0
#ADD . /go/src/myapp
#WORKDIR /go/src/myapp
#RUN go get -d -v ./...
#
## Install the package
#RUN go install -v ./...
#
##RUN go get myapp
##RUN go install
#EXPOSE 80
#
#ENTRYPOINT ["/go/bin/myapp"]

#FROM golang:alpine
#
#RUN apk add --no-cache git
#ENV CGO_ENABLED=0
#WORKDIR /app
##COPY . .
#COPY go.* .
#RUN go mod download
#
##RUN go get -d -v ./...
##
### Install the package
##RUN go install -v ./...
#
#COPY . .
#
#RUN cd . && go build -o main . && cp main ../../ && cd ../../
#
#CMD ["./main"]