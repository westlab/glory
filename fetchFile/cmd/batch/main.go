package main

import (
	"context"
	fetchFile "glory/fetchDrive"
	"log"
)

func main() {
	log.Print("start fetch files process")
	srv, err := fetchFile.GetCredential()
	if err != nil {
		log.Fatal(err)
	}

	if err = fetchFile.Handler(context.Background(), srv); err != nil {
		log.Fatal(err)
	}

}
