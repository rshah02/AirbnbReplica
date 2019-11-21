
FROM golang:latest

MAINTAINER Sindhu

# Expose port 8080 to the outside world
EXPOSE 8080

WORKDIR /app

COPY . .

#Getting required imports
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/sns
RUN go get github.com/gorilla/mux
RUN go get github.com/nu7hatch/gouuid
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/bson/primitive
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/mongo/options

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]