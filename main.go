package main

import (
	"fmt"

	"github.com/profsergiocosta/vm-translator/parser"
)

func main() {

	p := parser.New("StackTest.vm")
	fmt.Println(p.HasMoreCommands())

	for p.HasMoreCommands() {
		p.Advance()
		fmt.Println(p.CommandType())
	}

}
