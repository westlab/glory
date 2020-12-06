package main

import (
	"context"
	"log"

	"glory/countchars"
)

func main() {
	log.Print("start count process")

	err := countchars.Handler(context.Background())
	if err != nil {
		log.Fatal(err)
	}

}
