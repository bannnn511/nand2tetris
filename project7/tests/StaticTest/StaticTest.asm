//   
//   
//   
//   
//   
//   
// push constant 111
@111
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 333
@333
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 888
@888
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop static 8
@StaticTest.vm.8
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
// pop static 3
@StaticTest.vm.3
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
// pop static 1
@StaticTest.vm.1
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
// push static 3
@StaticTest.vm.3
D=M
@SP
A=M
M=D
@SP
M=M+1
// push static 1
@StaticTest.vm.1
D=M
@SP
A=M
M=D
@SP
M=M+1
// sub static 1
@SP
A=M-1
D=M
A=A-1
M=M-D
@SP
M=M-1
// push static 8
@StaticTest.vm.8
D=M
@SP
A=M
M=D
@SP
M=M+1
// add static 8
@SP
A=M-1
D=M
A=A-1
M=D+M
@SP
M=M-1
(END_PROGRAM)
@END_PROGRAM
0;JMP
