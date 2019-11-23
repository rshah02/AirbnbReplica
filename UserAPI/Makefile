
get-dependencies:
	rm -rf src/github.com
	go get -v github.com/gorilla/mux
	go get -v github.com/satori/go.uuid
	go get -v gopkg.in/mgo.v2
	go get -v gopkg.in/mgo.v2/bson
	go get firebase.google.com/go
	go get github.com/gorilla/handlers
	go get google.golang.org/api/option
	
build: get-dependencies
	  go build -o main .
	
run:
	./main