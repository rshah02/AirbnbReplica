package main
import (
    "fmt"
    "log"
    "net/http"
    
)

func main() {
    r := Router()
    fmt.Println("Starting server on the port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
