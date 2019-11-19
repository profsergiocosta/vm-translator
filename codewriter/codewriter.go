package codewriter

import (
	"fmt"
	"os"

	"github.com/profsergiocosta/vm-translator/command"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type CodeWriter struct {
	out *os.File
}

func New(pathName string) *CodeWriter {
	f, err := os.Create(pathName)
	check(err)

	code := &CodeWriter{out: f}

	return code
}

func (code *CodeWriter) WritePush(segment string, index int) {
	s := fmt.Sprintf("push %s %d\n", segment, index)
	code.out.WriteString(s)
}

func (code *CodeWriter) WritePop(segment string, index int) {
	s := fmt.Sprintf("pop %s %d\n", segment, index)
	code.out.WriteString(s)
}

func (code *CodeWriter) WriteArithmetic(cmd command.Command) {
	s := fmt.Sprintf("%s\n", cmd)
	code.out.WriteString(s)
}
func (code *CodeWriter) CloseFile() {
	code.out.Close()
}
