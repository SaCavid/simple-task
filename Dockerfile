FROM golang:1.15.8-alpine

# Copy the code into the container
COPY ./ ./
# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the application
RUN go build -o main .

# Export necessary port
EXPOSE 80

# Command to run when starting the container
CMD ["./main"]