@256
D=A
@SP
M=D
@300
D=A
@LCL
M=D
@400
D=A
@ARG
M=D
@6
D=A
@ARG
A=M
M=D
@3000
D=A
@ARG
A=M+1
M=D
//   
//   
//   
//   
//   
//   
//   
//   
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
// pop pointer 1
@1
D=A
@R3
A=D+A
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop that 0
@0
D=A
@THAT
A=D+M
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
// push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop that 1
@1
D=A
@THAT
A=D+M
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
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
// push constant 2
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
// sub constant 2
@SP
A=M-1
D=M
A=A-1
M=M-D
@SP
M=M-1
// pop argument 0
@0
D=A
@ARG
A=D+M
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
// pop argument 0
// label LOOP 0
(LOOP)
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
// if-goto COMPUTE_ELEMENT 0
@SP
A=M-1
D=M
A=A-1
@SP
M=M-1
@COMPUTE_ELEMENT
D;JNE
// goto END 0
@END
0;JMP
// goto END 0
// label COMPUTE_ELEMENT 0
(COMPUTE_ELEMENT)
// label COMPUTE_ELEMENT 0
// push that 0
@0
D=A
@THAT
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// push that 1
@1
D=A
@THAT
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// add that 1
@SP
A=M-1
D=M
A=A-1
M=D+M
@SP
M=M-1
// pop that 2
@2
D=A
@THAT
A=D+M
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
// pop that 2
// push pointer 1
@1
D=A
@R3
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
// push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
// add constant 1
@SP
A=M-1
D=M
A=A-1
M=D+M
@SP
M=M-1
// pop pointer 1
@1
D=A
@R3
A=D+A
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
// pop pointer 1
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
// push constant 1
@1
D=A
@SP
A=M
M=D
@SP
M=M+1
// sub constant 1
@SP
A=M-1
D=M
A=A-1
M=M-D
@SP
M=M-1
// pop argument 0
@0
D=A
@ARG
A=D+M
D=A
@R13
M=D
@SP
A=M-1
D=M
@R13
A=M
M=D
@SP
M=M-1
// goto LOOP 0
@LOOP
0;JMP
// goto LOOP 0
// label END 0
(END)
(END_PROGRAM)
@END_PROGRAM
0;JMP
