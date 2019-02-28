package parser

import (
	"io/ioutil"
	"regexp"
)

type Parser struct {
	tokens   []string
	position int
}

func New(fname string) *Parser {
	p := new(Parser)
	reComments, _ := regexp.Compile("//.*\n")
	reTokens, _ := regexp.Compile("[a-z][a-z]*|[1-9][0-9]*")
	code, _ := ioutil.ReadFile(fname)
	codeProc := reComments.ReplaceAllString(string(code), "")

	p.tokens = reTokens.FindAllString(codeProc, -1)
	p.position = 0

	return p
}

//métodos
func (self Parser) HasMoreCommands() bool { // publico tem que iniciar com maiuscula
	return self.position < len(self.tokens)
}
