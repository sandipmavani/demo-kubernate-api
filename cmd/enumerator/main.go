package main

import (
	"log"
	"net/http"
	"github.com/sandipmavani/demo-kubernate-api/handlers"
)

func main() {
	log.Println("Run the Enumerator")

	http.HandleFunc("/ping", handlers.Ping)
	
	http.HandleFunc("/enumerate", handlers.EnumerateRBAC)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
