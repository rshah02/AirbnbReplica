/*
	Gumball API in Go (Version 3)
	Uses MongoDB and RabbitMQ 
	(For use with Kong API Key)
*/

package main

import (
	"os"
    "fmt"
         
)

func main() {
           fmt.Printf("Welcome To ITzGeek\n")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	server := NewServer()
	server.Run(":" + port)
}
