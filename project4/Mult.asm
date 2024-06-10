// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
// The algorithm is based on repetitive addition.

@i
M=0 // i=0
@R2
M=0 // R2=0

(Loop)
@i
D=M // D=i
@R1
D=D-M // D=i-R1
@END
D;JEQ // if i=R1, goto END
@R0
D=M // D=R0
@R2
M=D+M // R2=i+R0
@i
M=M+1 // i=i+1
@Loop
0;JMP // goto Loop

(END)
@END
0;JMP // goto END