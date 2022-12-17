package main

import (
	"context"
	"fmt"

	"glory/fetchFile"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("target file(folder) missing")
	}
	target := os.Args[1]
	srv, err := fetchFile.GetCredential()
	if err != nil {
		log.Fatal(err)
	}

	output, err := fetchFile.FetchFileID(context.Background(), srv, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s id: %s\n", target, output)
}
