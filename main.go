package main

import (
	"fmt"
	"log"

	"github.com/Tulkdan/go-sse/internal"
)

func main() {
	fmt.Println("Starting server on 0.0.0.0:8000")

	if err := internal.NewServer("8000"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Killing server")
	}
}
