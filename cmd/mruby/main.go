package main

import (
	"log"

	"github.com/elct9620/mruby-go"
)

func main() {
	mrb, err := mruby.NewFromString("1 + 1")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(mrb.Header())
}
