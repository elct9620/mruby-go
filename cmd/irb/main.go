package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/elct9620/mruby-go"
)

const (
	exitCommand  = "exit"
	inputPrompt  = "> "
	outputPrompt = "=>"
)

func main() {
	mrb, err := mruby.New()
	if err != nil {
		panic(err)
	}

	fmt.Print("mirb - Embeded Interactive Ruby Shell (Go Edition)\n\n")
	fmt.Print(inputPrompt)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == exitCommand {
			break
		}

		ret, err := mrb.LoadString(line)

		// NOTE: ret may be error
		errRet, ok := ret.(error)
		if ok {
			err = errRet
		}

		if err != nil {
			fmt.Println(err)
		} else {
			if ret != nil {
				fmt.Printf("  %s %s\n", outputPrompt, mrb.Inspect(ret))
			}
		}

		fmt.Print(inputPrompt)
	}
}
