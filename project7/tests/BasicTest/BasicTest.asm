@10
D=A
@constant
A=M
M=D
@constant
M=M+1
LCL
D=M
0
D=D+A
@R13
M=D
LCL
AM=M-1
D=M
@R13
A=M
M=D
@21
D=A
@constant
A=M
M=D
@constant
M=M+1
@22
D=A
@constant
A=M
M=D
@constant
M=M+1
@36
D=A
@constant
A=M
M=D
@constant
M=M+1
THIS
D=M
6
D=D+A
@R13
M=D
THIS
AM=M-1
D=M
@R13
A=M
M=D
@42
D=A
@constant
A=M
M=D
@constant
M=M+1
@45
D=A
@constant
A=M
M=D
@constant
M=M+1
THAT
D=M
5
D=D+A
@R13
M=D
THAT
AM=M-1
D=M
@R13
A=M
M=D
THAT
D=M
2
D=D+A
@R13
M=D
THAT
AM=M-1
D=M
@R13
A=M
M=D
@510
D=A
@constant
A=M
M=D
@constant
M=M+1
R5
D=M
11
D=D+A
@R13
M=D
R5
AM=M-1
D=M
@R13
A=M
M=D
@0
D=A
@LCL
D=D+A
@LCL
A=M
M=D
@LCL
M=M+1
@5
D=A
@THAT
D=D+A
@THAT
A=M
M=D
@THAT
M=M+1

// add
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1

// sub
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
@6
D=A
@THIS
D=D+A
@THIS
A=M
M=D
@THIS
M=M+1
@6
D=A
@THIS
D=D+A
@THIS
A=M
M=D
@THIS
M=M+1

// add
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1

// sub
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
@11
D=A
@R5
D=D+A
@R5
A=M
M=D
@R5
M=M+1

// add
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
