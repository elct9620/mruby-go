package main

import (
	"bytes"
	"log"
	"os"

	"github.com/elct9620/mruby-go"
)

func main() {
	bin, err := mruby.Compile(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewBuffer(bin)
	rite, err := mruby.ReadRiteHeader(reader)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(rite)
}
