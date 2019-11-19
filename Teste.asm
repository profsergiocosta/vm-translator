@LCL // pop local 0
D=M
@0
D=D+A
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
A=M
M=D
@LCL // push local 10
D=M
@10
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
