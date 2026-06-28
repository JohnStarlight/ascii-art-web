package main

import (
	"fmt"
	"log"
	"net/http"

	"ascii-art-web/handlers"
)

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/ascii-art", handlers.AsciiArt)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
