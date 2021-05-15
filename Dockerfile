FROM golang:alpine

RUN apk add --no-cache git
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.* .
RUN go mod download

#RUN go get -d -v ./...
#
## Install the package
#RUN go install -v ./...

COPY . .

RUN cd . && go build -o main . && cp main ../../ && cd ../../

CMD ["./main"]