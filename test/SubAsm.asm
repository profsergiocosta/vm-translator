@100
M=1

@$RET1
D=A
@$Seta100$
0;JMP
($RET1)


@$RET2
D=A
@$Seta100$
0;JMP
($RET2)


// SUBROTINA
($Seta100$)
@R13    
M=D

@100
M=M+1

@R13
A=M
0;JMP
