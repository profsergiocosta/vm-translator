package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/profsergiocosta/vm-translator/codewriter"
	"github.com/profsergiocosta/vm-translator/command"
	"github.com/profsergiocosta/vm-translator/parser"
)

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		abs, _ := filepath.Abs(path)
		fmt.Printf("Could not find file or directory: %s \n", abs)
		os.Exit(1)
	}
	return fileInfo.IsDir()
}

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func translate(path string, code *codewriter.CodeWriter) {

	p := parser.New(path)
	code.SetFileName(path)
	//fmt.Println(p)
	//fmt.Println(path)
	for p.HasMoreCommands() {
		switch cmd := p.NextCommand().(type) {
		case command.Arithmetic:
			code.WriteArithmetic(cmd)
		case command.Push:
			code.WritePush(cmd.Segment, cmd.Index)
		case command.Pop:
			code.WritePop(cmd.Segment, cmd.Index)
		case command.Label:
			code.WriteLabel(cmd.Name)
		case command.Goto:
			code.WriteGoto(cmd.Label)
		case command.IFGoto:
			code.WriteIf(cmd.Label)
		case command.Return:
			code.WriteReturn()
		case command.CallFunction:
			code.WriteCall(cmd.FuncName, cmd.Args)
		case command.Function:
			code.WriteFunction(cmd.Name, cmd.Vars)
		}
	}

}

func main() {
	arg := os.Args[1:]

	if len(arg) == 1 {
		path := arg[0]

		if isDirectory(path) {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Fatal(err)
			}

			code := codewriter.New(filenameWithoutExtension(path) + "/" + filepath.Base(path) + ".asm")
			code.WriteInit()

			for _, f := range files {

				if filepath.Ext(f.Name()) == ".vm" {
					abs, _ := filepath.Abs(path + "/" + f.Name())
					fmt.Printf("Translating: %s \n", abs)
					fmt.Println(filenameWithoutExtension(path) + ".asm")
					translate(abs, code)

				}

			}
			code.CloseFile()

		} else {
			abs, _ := filepath.Abs(path)
			fmt.Printf("Translating: %s \n", abs)
			code := codewriter.New(filenameWithoutExtension(path) + ".asm")
			code.WriteInit()
			translate(path, code)
			code.WriteEnd()
			code.CloseFile()
		}

	}

}
