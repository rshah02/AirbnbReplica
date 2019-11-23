
FROM golang:latest

MAINTAINER Mrinalini

# Expose port 3000 to the outside world
EXPOSE 3000

WORKDIR /app

COPY main /app 

#COPY . .

#Getting required imports
#RUN go get github.com/gorilla/mux
#RUN go get github.com/satori/go.uuid
#RUN go get gopkg.in/mgo.v2
#RUN go get gopkg.in/mgo.v2/bson
#RUN go get firebase.google.com/go

# Build the Go app
#RUN go build -o main .

# Command to run the executable
CMD ["./main"]