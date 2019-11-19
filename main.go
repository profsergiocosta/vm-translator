package main

import (
	"github.com/profsergiocosta/vm-translator/codewriter"
	"github.com/profsergiocosta/vm-translator/command"
	"github.com/profsergiocosta/vm-translator/parser"
)

func main() {

	p := parser.New("StackTest.vm")
	code := codewriter.New("saida.asm")

	for p.HasMoreCommands() {

		switch cmd := p.NextCommand().(type) {
		case command.Arithmetic:
			code.WriteArithmetic(cmd)
		case command.Push:
			code.WritePush(cmd.Segment, cmd.Index)
		case command.Pop:
			code.WritePop(cmd.Segment, cmd.Index)
		}
	}
	code.CloseFile()

}
