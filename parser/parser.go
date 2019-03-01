package parser

import (
	"io/ioutil"
	"regexp"
)

type Parser struct {
	tokens    []string
	position  int
	currToken string
}

func New(fname string) *Parser {
	p := new(Parser)
	reComments, err := regexp.Compile("//.*\n")
	if err != nil {
		// tratar o erro aqui
		panic("Error")
	}
	reTokens, _ := regexp.Compile("[a-z][a-z]*|[1-9][0-9]*")
	code, _ := ioutil.ReadFile(fname)
	codeProc := reComments.ReplaceAllString(string(code), "")

	p.tokens = reTokens.FindAllString(codeProc, -1)
	p.position = 0

	return p
}

//métodos
func (self *Parser) HasMoreCommands() bool { // publico tem que iniciar com maiuscula
	return self.position < len(self.tokens)
}

func (self *Parser) Advance() {
	self.currToken = self.tokens[self.position]
	self.position++
}

func (self *Parser) CommandType() string {
	switch self.currToken {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return "arithmetic"
	default:
		return self.currToken
	}
}

func (self *Parser) Arg1() string {

	if self.CommandType() == "arithmetic" {
		return self.currToken
	} else {
		self.Advance()
		return self.currToken
	}

}
func (self *Parser) Arg2() string {
	self.Advance()
	return self.currToken
}
