//   
//   
//   
//   
//   
//   
//   
// function SimpleFunction.test 2
(SimpleFunction.test)
@SP
A=M
M=0
@SP
M=M+1
@SP
A=M
M=0
@SP
M=M+1
// push local 0
@0
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// push local 1
@1
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// add local 1
@SP
A=M-1
D=M
A=A-1
M=D+M
@SP
M=M-1
// not local 1
@SP
A=M-1
M=!M
// push argument 0
@0
D=A
@ARG
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// add argument 0
@SP
A=M-1
D=M
A=A-1
M=D+M
@SP
M=M-1
// push argument 1
@1
D=A
@ARG
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// sub argument 1
@SP
A=M-1
D=M
A=A-1
M=M-D
@SP
M=M-1
// return argument 1
@LCL
D=M
@R13
M=D
@5
A=D-A
M=D
@R14
M=D
@SP
A=M-1
D=M
A=A-1
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@R13
AM=M-1
D=M
@THAT
M=D
@R13
AM=M-1
D=M
@THIS
M=D
@R13
AM=M-1
D=M
@ARG
M=D
@R13
AM=M-1
D=M
@LCL
M=D
@R14
A=M
0;JMP
(END_PROGRAM)
@END_PROGRAM
0;JMP
