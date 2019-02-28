package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {

	reComments, _ := regexp.Compile("//.*\n")

	reTokens, _ := regexp.Compile("[a-z][a-z]*|[1-9][0-9]*")

	code, _ := ioutil.ReadFile("StackTest.vm")

	codeProc := reComments.ReplaceAllString(string(code), "")

	tokens := reTokens.FindAllString(codeProc, -1)
	fmt.Printf("%T\n", tokens)

}
