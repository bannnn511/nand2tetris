// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input. 
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, 
// the screen should be cleared.

@offset // set offset
M=0 // place at upper left 

(KBDCheck)
@KBD
D=M // D = RAM[24576] -> key pressed 
@White
D;JEQ // if D == 0, goto White
@Black
0;JMP // if D >= 0, goto Black

(Black)
@place
D=M
@8192 // max screen size
D=D-A
@KBDCheck
D;JGE
@place
D=M
@SCREEN
A=D+A // screen + offset
M=-1
@place
M=M+1
@KBDCheck
0;JMP

(White)
@place
D=M
@SCREEN
A=D+A // screen + offset
M=0
@place
M=M-1
@KBDCheck
0;JMP


@KBDCheck
0;JMP