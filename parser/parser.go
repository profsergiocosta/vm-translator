package parser

import (
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/profsergiocosta/vm-translator/command"
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
	reTokens, _ := regexp.Compile("[a-z][a-z]*|0|[1-9][0-9]*")
	code, _ := ioutil.ReadFile(fname)
	codeProc := reComments.ReplaceAllString(string(code), "")

	p.tokens = reTokens.FindAllString(codeProc, -1)
	p.position = 0
	//p.Advance()
	return p
}

func (p *Parser) HasMoreCommands() bool {
	return p.position < len(p.tokens)
}

func (p *Parser) Advance() {
	p.currToken = p.tokens[p.position]
	p.position++
}

func (p *Parser) NextCommand() command.Command {
	p.Advance()

	switch p.currToken {

	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return command.Arithmetic{Name: p.currToken}
	case "label", "if-goto", "goto":
		cmd := p.currToken

		p.Advance()
		arg1 := p.currToken
		switch cmd {
		case "label":
			return command.Label{Name: arg1}
		case "goto":
			return command.Goto{Label: arg1}
		case "if-goto":
			return command.IFGoto{Label: arg1}

		}
	case "push", "pop", "function", "call":
		cmd := p.currToken
		p.Advance()
		arg1 := p.currToken
		p.Advance()
		arg2, _ := strconv.Atoi(p.currToken)
		switch cmd {
		case "push":
			return command.Push{Segment: arg1, Index: arg2}
		case "pop":
			return command.Pop{Segment: arg1, Index: arg2}
		}

	}

	return command.UndefinedCommand{}
}
