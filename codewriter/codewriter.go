package codewriter

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/profsergiocosta/vm-translator/command"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

type CodeWriter struct {
	out        *os.File
	moduleName string
}

func New(pathName string) *CodeWriter {
	f, err := os.Create(pathName)
	check(err)

	code := &CodeWriter{out: f}

	return code
}

func (code *CodeWriter) write(s string) {
	code.out.WriteString(fmt.Sprintf("%s\n", s))

}

func (code *CodeWriter) segmentPointer(segment string, index int) string {
	switch segment {
	case "local":
		return "LCL"
	case "argument":
		return "ARG"
	case "this", "that":
		return strings.ToUpper(segment)
	case "temp":
		return fmt.Sprintf("R%d", 5+index)
	case "pointer":
		return fmt.Sprintf("R%d", 3+index)
	case "static":
		return fmt.Sprintf("%s.%s", code.moduleName, index)
	default:
		return "ERROR"
	}

}

func (code *CodeWriter) SetFileName(pathName string) {
	code.moduleName = filenameWithoutExtension(path.Base(pathName))
}

func (code *CodeWriter) WritePush(seg string, index int) {
	switch seg {
	case "constant":
		code.write(fmt.Sprintf("@%d // push %s %d", index, seg, index))
		code.write("D=A")
		code.write("@SP")
		code.write("A=M")
		code.write("M=D")
		code.write("@SP")
		code.write("M=M+1")
	case "static", "temp", "pointer":
		code.write(fmt.Sprintf("@%s // push %s %d", code.segmentPointer(seg, index), seg, index))
		code.write("D=M")
		code.write("@SP")
		code.write("A=M")
		code.write("M=D")
		code.write("@SP")
		code.write("M=M+1")
	case "local", "argument", "this", "that":
		code.write(fmt.Sprintf("@%s // push %s %d", code.segmentPointer(seg, index), seg, index))
		code.write("D=M")
		code.write(fmt.Sprintf("@%d", index))
		code.write("A=D+A")
		code.write("D=M")
		code.write("@SP")
		code.write("A=M")
		code.write("M=D")
		code.write("@SP")
		code.write("M=M+1")
	default:

	}
}

func (code *CodeWriter) WritePop(seg string, index int) {
	switch seg {
	case "static", "temp", "pointer":
		code.write(fmt.Sprintf("@SP // pop %s %d", seg, index))
		code.write("M=M-1")
		code.write("A=M")
		code.write("D=M")
		code.write(fmt.Sprintf("@%s", code.segmentPointer(seg, index)))
		code.write("M=D")
	case "local", "argument", "this", "that":
		code.write(fmt.Sprintf("@%s // pop %s %d", code.segmentPointer(seg, index), seg, index))
		code.write("D=M")
		code.write(fmt.Sprintf("@%d", index))
		code.write("D=D+A")
		code.write("@R13")
		code.write("M=D")
		code.write("@SP")
		code.write("M=M-1")
		code.write("A=M")
		code.write("D=M")
		code.write("@R13")
		code.write("A=M")
		code.write("M=D")
	default:

	}

}

func (code *CodeWriter) WriteArithmetic(cmd command.Arithmetic) {
	switch cmd.Name {
	case "add":
		code.writeArithmeticAdd()
	case "sub":
		code.writeArithmeticSub()
	}
}

func (code *CodeWriter) writeArithmeticAdd() {
	code.write("@SP // add")
	code.write("M=M-1")
	code.write("A=M")
	code.write("D=M")
	code.write("A=A-1")
	code.write("M=D+M")
}

func (code *CodeWriter) writeArithmeticSub() {
	code.write("@SP // sub")
	code.write("M=M-1")
	code.write("A=M")
	code.write("D=M")
	code.write("A=A-1")
	code.write("M=M-D")
}

func (code *CodeWriter) CloseFile() {
	code.out.Close()
}
