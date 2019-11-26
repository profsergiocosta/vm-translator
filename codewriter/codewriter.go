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
	funcName   string
	labelCount int
	callCount  int
}

func New(pathName string) *CodeWriter {
	f, err := os.Create(pathName)
	check(err)

	code := &CodeWriter{out: f}
	code.labelCount = 0
	code.callCount = 0

	code.funcName = ""

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
		return fmt.Sprintf("%s.%d", code.moduleName, index)
	default:
		return "ERROR"
	}

}

func (code *CodeWriter) WriteInit() {
	code.write("@256")
	code.write("D=A")
	code.write("@SP")
	code.write("M=D")
	code.WriteCall("Sys.init", 0)
}

func (code *CodeWriter) WriteEnd() {
	code.write("(LOOP_INFINITO)")
	code.write("@LOOP_INFINITO")
	code.write("0;JMP")

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
	case "neg":
		code.writeArithmeticNeg()
	case "eq":
		code.writeArithmeticEq()
	case "gt":
		code.writeArithmeticGt()
	case "lt":
		code.writeArithmeticLt()
	case "and":
		code.writeArithmeticAnd()
	case "or":
		code.writeArithmeticOr()
	case "not":
		code.writeArithmeticNot()
	default:
	}
}

func (code *CodeWriter) writeBinaryArithmetic() {
	code.write("@SP")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("A=A-1")
}

func (code *CodeWriter) writeArithmeticAdd() {
	code.writeBinaryArithmetic()
	code.write("M=D+M")
}

func (code *CodeWriter) writeArithmeticSub() {
	code.writeBinaryArithmetic()
	code.write("M=M-D")
}

func (code *CodeWriter) writeArithmeticAnd() {
	code.writeBinaryArithmetic()
	code.write("M=D&M")
}

func (code *CodeWriter) writeArithmeticOr() {
	code.writeBinaryArithmetic()
	code.write("M=D|M")
}

func (code *CodeWriter) writeUnaryArithmetic() {
	code.write("@SP")
	code.write("A=M")
	code.write("A=A-1")
}

func (code *CodeWriter) writeArithmeticNeg() {
	code.writeUnaryArithmetic()
	code.write("M=-M")
}

func (code *CodeWriter) writeArithmeticNot() {
	code.writeUnaryArithmetic()
	code.write("M=!M")
}

func (code *CodeWriter) writeArithmeticEq() {

	label := fmt.Sprintf("JEQ_%s_%d", code.moduleName, code.labelCount)
	code.write("@SP // eq")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@SP")
	code.write("AM=M-1")
	code.write("D=M-D")
	code.write("@" + label)
	code.write("D;JEQ")
	code.write("D=1")
	code.write("(" + label + ")")
	code.write("D=D-1")
	code.write("@SP")
	code.write("A=M")
	code.write("M=D")
	code.write("@SP")
	code.write("M=M+1")

	code.labelCount++
}

func (code *CodeWriter) writeArithmeticGt() {

	labelTrue := fmt.Sprintf("JGT_TRUE_%s_%d", code.moduleName, code.labelCount)
	labelFalse := fmt.Sprintf("JGT_FALSE_%s_%d", code.moduleName, code.labelCount)

	code.write("@SP // gt")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@SP")
	code.write("AM=M-1")
	code.write("D=M-D")
	code.write("@" + labelTrue)
	code.write("D;JGT")
	code.write("D=0")
	code.write("@" + labelFalse)
	code.write("0;JMP")
	code.write("(" + labelTrue + ")")
	code.write("D=-1")
	code.write("(" + labelFalse + ")")
	code.write("@SP")
	code.write("A=M")
	code.write("M=D")
	code.write("@SP")
	code.write("M=M+1")

	code.labelCount++
}

func (code *CodeWriter) writeArithmeticLt() {

	labelTrue := fmt.Sprintf("JLT_TRUE_%s_%d", code.moduleName, code.labelCount)
	labelFalse := fmt.Sprintf("JLT_FALSE_%s_%d", code.moduleName, code.labelCount)

	code.write("@SP // lt")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@SP")
	code.write("AM=M-1")
	code.write("D=M-D")
	code.write("@" + labelTrue + "")
	code.write("D;JLT")
	code.write("D=0")
	code.write("@" + labelFalse + "")
	code.write("0;JMP")
	code.write("(" + labelTrue + ")")
	code.write("D=-1")
	code.write("(" + labelFalse + ")")
	code.write("@SP")
	code.write("A=M")
	code.write("M=D")
	code.write("@SP")
	code.write("M=M+1")

	code.labelCount++
}

func (code *CodeWriter) WriteLabel(label string) {
	code.write("(" + label + ")")
}

func (code *CodeWriter) WriteGoto(label string) {
	code.write("@" + label)
	code.write("0;JMP")
}

func (code *CodeWriter) WriteIf(label string) {
	code.write("@SP")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("M=0")
	code.write("@" + label)
	code.write("D;JNE")

}

func (code *CodeWriter) WriteFunction(funcName string, nLocals int) {

	loopLabel := funcName + "_INIT_LOCALS_LOOP"
	loopEndLabel := funcName + "_INIT_LOCALS_END"

	code.funcName = funcName

	code.write("(" + funcName + ")" + "// initializa local variables")
	code.write(fmt.Sprintf("@%d", nLocals))
	code.write("D=A")
	code.write("@R13") // temp
	code.write("M=D")
	code.write("(" + loopLabel + ")")
	code.write("@" + loopEndLabel)
	code.write("D;JEQ")
	code.write("@0")
	code.write("D=A")
	code.write("@SP")
	code.write("A=M")
	code.write("M=D")
	code.write("@SP")
	code.write("M=M+1")
	code.write("@R13")
	code.write("MD=M-1")
	code.write("@" + loopLabel)
	code.write("0;JMP")
	code.write("(" + loopEndLabel + ")")

}

func (code *CodeWriter) writeFramePush(value string) {
	code.write("@" + value)
	code.write("D=M")
	code.write("@SP")
	code.write("A=M")
	code.write("M=D")
	code.write("@SP")
	code.write("M=M+1")
}

func (code *CodeWriter) WriteCall(funcName string, numArgs int) {

	/*
	   push return-address     // (using the label declared below)
	   push LCL                // save LCL of the calling function
	   push ARG                // save ARG of the calling function
	   push THIS               // save THIS of the calling function
	   push THAT               // save THAT of the calling function
	   ARG = SP-n-5            // reposition ARG (n = number of args)
	   LCL = SP                // reposiiton LCL
	   goto f                  // transfer control
	   (return-address)        // declare a label for the return-address
	*/

	comment := fmt.Sprintf("// call %s %d", funcName, numArgs)

	returnAddr := fmt.Sprintf("%s_RETURN_%d", funcName, code.callCount)
	code.callCount++

	code.write(fmt.Sprintf("@%s %s", returnAddr, comment)) // push return-addr
	code.write("D=A")
	code.write("@SP")
	code.write("A=M")
	code.write("M=D")
	code.write("@SP")
	code.write("M=M+1")

	code.writeFramePush("LCL")
	code.writeFramePush("ARG")
	code.writeFramePush("THIS")
	code.writeFramePush("THAT")

	code.write(fmt.Sprintf("@%d", numArgs)) // ARG = SP-n-5
	code.write("D=A")
	code.write("@5")
	code.write("D=D+A")
	code.write("@SP")
	code.write("D=M-D")
	code.write("@ARG")
	code.write("M=D")

	code.write("@SP") // LCL = SP
	code.write("D=M")
	code.write("@LCL")
	code.write("M=D")

	code.WriteGoto(funcName)

	code.write("(" + returnAddr + ")") // (return-address)

}

func (code *CodeWriter) WriteReturn() {

	/*
	   FRAME = LCL         // FRAME is a temporary var
	   RET = *(FRAME-5)    // put the return-address in a temporary var
	   *ARG = pop()        // reposition the return value for the caller
	   SP = ARG + 1        // restore SP of the caller
	   THAT = *(FRAME - 1) // restore THAT of the caller
	   THIS = *(FRAME - 2) // restore THIS of the caller
	   ARG = *(FRAME - 3)  // restore ARG of the caller
	   LCL = *(FRAME - 4)  // restore LCL of the caller
	   goto RET            // goto return-address (in the caller's code)
	*/

	code.write("@LCL") // FRAME = LCL
	code.write("D=M")

	code.write("@R13") // R13 -> FRAME
	code.write("M=D")

	code.write("@5") // RET = *(FRAME-5)
	code.write("A=D-A")
	code.write("D=M")
	code.write("@R14") // R14 -> RET
	code.write("M=D")

	code.write("@SP") // *ARG = pop()
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@ARG")
	code.write("A=M")
	code.write("M=D")

	code.write("D=A") // SP = ARG+1
	code.write("@SP")
	code.write("M=D+1")

	code.write("@R13") // THAT = *(FRAME-1)
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@THAT")
	code.write("M=D")

	code.write("@R13") // THIS = *(FRAME-2)
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@THIS")
	code.write("M=D")

	code.write("@R13") // ARG = *(FRAME-3)
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@ARG")
	code.write("M=D")

	code.write("@R13") // LCL = *(FRAME-4)
	code.write("AM=M-1")
	code.write("D=M")
	code.write("@LCL")
	code.write("M=D")

	code.write("@R14") // goto RET
	code.write("A=M")
	code.write("0;JMP")

}

func (code *CodeWriter) CloseFile() {
	code.out.Close()
}
