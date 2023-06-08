package main

import (
	"log"
	"os"

	"github.com/elct9620/mruby-go"
)

func main() {
	bin, err := mruby.Compile(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Compiled binary: ", bin)
}
