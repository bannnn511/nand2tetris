@111
D=A
@constant
A=M
M=D
@constant
M=M+1
@333
D=A
@constant
A=M
M=D
@constant
M=M+1
@888
D=A
@constant
A=M
M=D
@constant
M=M+1
@R24
D=A
@R13
M=D
@R24
AM=M-1
D=M
@R13
A=M
M=D
@R19
D=A
@R13
M=D
@R19
AM=M-1
D=M
@R13
A=M
M=D
@R17
D=A
@R13
M=D
@R17
AM=M-1
D=M
@R13
A=M
M=D
@R19
D=A
@R19
A=M
M=D
@R19
M=M+1
@R17
D=A
@R17
A=M
M=D
@R17
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
@R24
D=A
@R24
A=M
M=D
@R24
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