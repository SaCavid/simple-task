FROM golang:alpine

RUN apk add --no-cache git
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .

RUN cd . && go build -o main . && cp main ../../ && cd ../../

CMD ["./main"]