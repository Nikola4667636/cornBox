package main

import (
	"flag"
	"fmt"
	"log"

	"cronBox/secretstore"
)

func main() {
	plain := flag.String("value", "", "Secret value to encrypt")
	flag.Parse()

	encrypted, err := secretstore.Encrypt(*plain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encrypted)
}
