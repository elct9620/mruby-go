package main

import (
	"bytes"
	"log"
	"os"

	"github.com/elct9620/mruby-go"
)

func main() {
	mrb, err := mruby.New()
	if err != nil {
		log.Fatal(err)
	}

	compiled, err := mruby.Compile(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	res, err := mrb.LoadIRep(bytes.NewBuffer(compiled))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
