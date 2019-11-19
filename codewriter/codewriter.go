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
	labelCount int
}

func New(pathName string) *CodeWriter {
	f, err := os.Create(pathName)
	check(err)

	code := &CodeWriter{out: f}
	code.labelCount = 0

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

func (code *CodeWriter) writeArithmeticNeg() {
	code.write("@SP // neg")
	code.write("A=M")
	code.write("A=A-1")
	code.write("M=-M")
}

func (code *CodeWriter) writeArithmeticAnd() {
	code.write("@SP // and")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("A=A-1")
	code.write("M=D&M")
}

func (code *CodeWriter) writeArithmeticOr() {
	code.write("@SP // or")
	code.write("AM=M-1")
	code.write("D=M")
	code.write("A=A-1")
	code.write("M=D|M")
}

func (code *CodeWriter) writeArithmeticNot() {
	code.write("@SP // not")
	code.write("A=M")
	code.write("A=A-1")
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
	code.write("@" + labelTrue)
	code.write("D;JLT")
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

func (code *CodeWriter) writeLabel (label string ) {
	code.write ("("+label+")");
}

func (code *CodeWriter) writeGoto (label string) {
	code.write ("@" + label );
	code.write ("0;JMP");
}

func (code *CodeWriter) writeIf (label string) {
	code.write ("@SP");
	code.write ("AM=M-1");
	code.write ("D=M");
	code.write ("M=0");
	code.write ("@"+label);
	code.write ("D;JNE");
}

func (code *CodeWriter) writeFunction(funcName string,nLocals int){


	loopLabel := funcName + "_INIT_LOCALS_LOOP"
	loopEndLabel := funcName + "_INIT_LOCALS_END"

    code.write("(" + funcName + ")" + "// initializa local variables");
	code.write( fmt.Sprintf("@%d", nLocals));
	   code.write("D=A");
	    code.write("@R13"); // temp
	    code.write("M=D");
	    code.write("(" + loopLabel + ")");
	    code.write("@" + loopEndLabel);
	    code.write("D;JEQ");
	    code.write("@0");
	    code.write("D=A");
	    code.write("@SP");
	    code.write("A=M");
	    code.write("M=D");
	    code.write("@SP");
	    code.write("M=M+1");
	    code.write("@R13");
	    code.write("MD=M-1");
	    code.write("@" + loopLabel);
	    code.write("0;JMP");
	    code.write("(" + loopEndLabel + ")");               
}


void CodeWriter::writeReturn () {

	
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








    write("@LCL"); // FRAME = LCL

    write("D=M");

    write("@R13"); void CodeWriter::writeReturn () {
	
	/*
	write ("@LCL // return "); // endFrame = LCL
	write ("D=M");
	write ("@R13");
	write ("M=D");
	
	write ("@5");  // retAdress = *(endFrame - 5) // <<"@5\nD=A\n@R13\nA=M-D\nD=M\n@R14\nM=D\n"
	write ("A=D-A");
	write ("D=M");
	write ("@R14");
	write ("M=D");
	
	write ("@SP // *arg = pop ()"); // *arg = pop ()
	write ("AM=M-1");
	write ("D=M");
	write ("@ARG");
	write ("A=M");
	write ("M=D");
	
	write ("D=A // Sp = arg+ 1"); 
	write ("@SP");
	write ("M=D+1");
	
	write("@R13"); // THAT = *(FRAME-1)
	write("AM=M-1");
	write("D=M");
	write("@THAT");
	write("M=D");
	
	write("@R13"); // THIS = *(FRAME-4)
	write("AM=M-1");
	write("D=M");
	write("@THIS");
	write("M=D");
	
	write("@R13"); // ARG = *(FRAME-3)
	write("AM=M-1");
	write("D=M");
	write("@ARG");
	write("M=D");
	
	write("@R13"); // THAT = *(FRAME-4)
	write("AM=M-1");
	write("D=M");
	write("@LCL");
	write("M=D");
	
	
	write ("@R14"); // goto ret
	write ("A=M");
	write ("0;JMP");
	*/
	
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








    write("@LCL"); // FRAME = LCL

    write("D=M");

    write("@R13"); // R13 -> FRAME

    write("M=D");




    write("@5"); // RET = *(FRAME-5)

    write("A=D-A");

    write("D=M");

    write("@R14"); // R14 -> RET

    write("M=D");




    write("@SP"); // *ARG = pop()

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("A=M");

    write("M=D");




    write("D=A"); // SP = ARG+1

    write("@SP");

    write("M=D+1");




    write("@R13"); // THAT = *(FRAME-1)

    write("AM=M-1");

    write("D=M");

    write("@THAT");

    write("M=D");




    write("@R13"); // THIS = *(FRAME-2)

    write("AM=M-1");

    write("D=M");

    write("@THIS");

    write("M=D");




    write("@R13"); // ARG = *(FRAME-3)

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("M=D");




    write("@R13"); // LCL = *(FRAME-4)

    write("AM=M-1");

    write("D=M");

    write("@LCL");

    write("M=D");




    write("@R14"); // goto RET

    write("A=M");

    write("0;JMP");
  

}

    write("M=D");




    write("@5"); //void CodeWriter::writeReturn () {
	
	/*
	write ("@LCL // return "); // endFrame = LCL
	write ("D=M");
	write ("@R13");
	write ("M=D");
	
	write ("@5");  // retAdress = *(endFrame - 5) // <<"@5\nD=A\n@R13\nA=M-D\nD=M\n@R14\nM=D\n"
	write ("A=D-A");
	write ("D=M");
	write ("@R14");
	write ("M=D");
	
	write ("@SP // *arg = pop ()"); // *arg = pop ()
	write ("AM=M-1");
	write ("D=M");
	write ("@ARG");
	write ("A=M");
	write ("M=D");
	
	write ("D=A // Sp = arg+ 1"); 
	write ("@SP");
	write ("M=D+1");
	
	write("@R13"); // THAT = *(FRAME-1)
	write("AM=M-1");
	write("D=M");
	write("@THAT");
	write("M=D");
	
	write("@R13"); // THIS = *(FRAME-4)
	write("AM=M-1");
	write("D=M");
	write("@THIS");
	write("M=D");
	
	write("@R13"); // ARG = *(FRAME-3)
	write("AM=M-1");
	write("D=M");
	write("@ARG");
	write("M=D");
	
	write("@R13"); // THAT = *(FRAME-4)
	write("AM=M-1");
	write("D=M");
	write("@LCL");
	write("M=D");
	
	
	write ("@R14"); // goto ret
	write ("A=M");
	write ("0;JMP");
	*/
	
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








    write("@LCL"); // FRAME = LCL

    write("D=M");

    write("@R13"); // R13 -> FRAME

    write("M=D");




    write("@5"); // RET = *(FRAME-5)

    write("A=D-A");

    write("D=M");

    write("@R14"); // R14 -> RET

    write("M=D");




    write("@SP"); // *ARG = pop()

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("A=M");

    write("M=D");




    write("D=A"); // SP = ARG+1

    write("@SP");

    write("M=D+1");




    write("@R13"); // THAT = *(FRAME-1)

    write("AM=M-1");

    write("D=M");

    write("@THAT");

    write("M=D");




    write("@R13"); // THIS = *(FRAME-2)

    write("AM=M-1");

    write("D=M");

    write("@THIS");

    write("M=D");




    write("@R13"); // ARG = *(FRAME-3)

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("M=D");




    write("@R13"); // LCL = *(FRAME-4)

    write("AM=M-1");

    write("D=M");

    write("@LCL");

    write("M=D");




    write("@R14"); // goto RET

    write("A=M");

    write("0;JMP");
  

}

    write("A=D-A");void CodeWriter::writeReturn () {
	
	/*
	write ("@LCL // return "); // endFrame = LCL
	write ("D=M");
	write ("@R13");
	write ("M=D");
	
	write ("@5");  // retAdress = *(endFrame - 5) // <<"@5\nD=A\n@R13\nA=M-D\nD=M\n@R14\nM=D\n"
	write ("A=D-A");
	write ("D=M");
	write ("@R14");
	write ("M=D");
	
	write ("@SP // *arg = pop ()"); // *arg = pop ()
	write ("AM=M-1");
	write ("D=M");
	write ("@ARG");
	write ("A=M");
	write ("M=D");
	
	write ("D=A // Sp = arg+ 1"); 
	write ("@SP");
	write ("M=D+1");
	
	write("@R13"); // THAT = *(FRAME-1)
	write("AM=M-1");
	write("D=M");
	write("@THAT");
	write("M=D");
	
	write("@R13"); // THIS = *(FRAME-4)
	write("AM=M-1");
	write("D=M");
	write("@THIS");
	write("M=D");
	
	write("@R13"); // ARG = *(FRAME-3)
	write("AM=M-1");
	write("D=M");
	write("@ARG");
	write("M=D");
	
	write("@R13"); // THAT = *(FRAME-4)
	write("AM=M-1");
	write("D=M");
	write("@LCL");
	write("M=D");
	
	
	write ("@R14"); // goto ret
	write ("A=M");
	write ("0;JMP");
	*/
	
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








    write("@LCL"); // FRAME = LCL

    write("D=M");

    write("@R13"); // R13 -> FRAME

    write("M=D");




    write("@5"); // RET = *(FRAME-5)

    write("A=D-A");

    write("D=M");

    write("@R14"); // R14 -> RET

    write("M=D");




    write("@SP"); // *ARG = pop()

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("A=M");

    write("M=D");




    write("D=A"); // SP = ARG+1

    write("@SP");

    write("M=D+1");




    write("@R13"); // THAT = *(FRAME-1)

    write("AM=M-1");

    write("D=M");

    write("@THAT");

    write("M=D");




    write("@R13"); // THIS = *(FRAME-2)

    write("AM=M-1");

    write("D=M");

    write("@THIS");

    write("M=D");




    write("@R13"); // ARG = *(FRAME-3)

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("M=D");




    write("@R13"); // LCL = *(FRAME-4)

    write("AM=M-1");

    write("D=M");

    write("@LCL");

    write("M=D");




    write("@R14"); // goto RET

    write("A=M");

    write("0;JMP");
  

}

    write("D=M");

    write("@R14"); // R14 -> RET

    write("M=D");




    write("@SP"); // *ARG = pop()

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("A=M");

    write("M=D");




    write("D=A"); // SP = ARG+1

    write("@SP");

    write("M=D+1");




    write("@R13"); // THAT = *(FRAME-1)

    write("AM=M-1");

    write("D=M");

    write("@THAT");

    write("M=D");




    write("@R13"); // THIS = *(FRAME-2)

    write("AM=M-1");

    write("D=M");

    write("@THIS");

    write("M=D");




    write("@R13"); // ARG = *(FRAME-3)

    write("AM=M-1");

    write("D=M");

    write("@ARG");

    write("M=D");




    write("@R13"); // LCL = *(FRAME-4)

    write("AM=M-1");

    write("D=M");

    write("@LCL");

    write("M=D");




    write("@R14"); // goto RET

    write("A=M");

    write("0;JMP");
  

}

func (code *CodeWriter) CloseFile() {
	code.out.Close()
}
