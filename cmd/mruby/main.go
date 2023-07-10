package main

import (
	"bytes"
	"log"
	"os"

	"github.com/elct9620/mruby-go"
)

func main() {
	mrb := mruby.New()

	compiled, err := mruby.Compile(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	res, err := mrb.Load(bytes.NewBuffer(compiled))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.(string))
}
