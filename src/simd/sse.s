#include "textflag.h"

TEXT ·Add(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	ADDPS     X1, X0
	MOVUPS    X0, ret+32(FP)
	RET 

TEXT ·Sub(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	SUBPS     X1, X0
	MOVUPS    X0, ret+32(FP)
	RET

TEXT ·Min(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	MINPS     X1, X0
	MOVUPS    X0, ret+32(FP)
	RET

TEXT ·Max(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	MAXPS     X1, X0
	MOVUPS    X0, ret+32(FP)
	RET
	
// this one is kinda dumb..
TEXT ·Dot(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	MULPS     X1, X0
	HADDPS    X0, X0
	HADDPS    X0, X0
	MOVUPS    X0, ret+32(FP)
	RET

TEXT ·Cross(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	MOVAPS    X0, X2
	MOVAPS    X1, X3

	SHUFPS    $0xd8, X0, X0
	SHUFPS    $0xe1, X1, X1
	MULPS     X1, X0

	SHUFPS    $0xe1, X2, X2
	SHUFPS    $0xd8, X3, X3
	MULPS     X3, X2

	SUBPS     X2, X0	
	MOVUPS    X0, ret+32(FP)
	RET

// single lane box intersect: t0 and t1 compared outside of this func
TEXT ·BoxIntersect(SB), NOSPLIT, $0-64
	MOVUPS    origins+0(FP), X0
	MOVUPS    directions+16(FP), X1
	MOVUPS    mins+32(FP), X2
	MOVUPS    maxs+48(FP), X3
	RCPPS     X1, X1 
	SUBPS     X0, X2
	MULPS     X1, X2
	SUBPS     X0, X3
	MULPS     X1, X3
	MOVAPS    X2, X0
	MINPS     X3, X2
	MAXPS     X0, X3
	MOVUPS    X2, ret+64(FP)
	MOVUPS    X3, ret+80(FP)
	RET

// t1s is a quadfloat set to max float values
DATA t1s<>+0x00(SB)/4, $0xffffffff
DATA t1s<>+0x04(SB)/4, $0xffffffff
DATA t1s<>+0x08(SB)/4, $0xffffffff
DATA t1s<>+0x0c(SB)/4, $0xffffffff
GLOBL t1s<>(SB), (NOPTR+RODATA), $16

// 4 lane box intersect: output is t0 values for the 4 streams
// t0 is set to 0 if no hit was found
TEXT ·Box4Intersect(SB), NOSPLIT, $0-192
	MOVUPS    o4x+0(FP), X0
	MOVUPS    o4y+16(FP), X1
	MOVUPS    o4z+32(FP), X2
	MOVUPS    d4x+48(FP), X3
	MOVUPS    d4y+64(FP), X4
	MOVUPS    d4z+80(FP), X5
	MOVUPS    min4x+96(FP), X6
	MOVUPS    min4y+112(FP), X7
	MOVUPS    min4z+128(FP), X8
	MOVUPS    max4x+144(FP), X9
	MOVUPS    max4y+160(FP), X10
	MOVUPS    max4z+176(FP), X11
	MOVAPS    t1s<>(SB), X12

	// these 3 blocks each create a mask, one for each dimension
	// in the meantime t0 is constructed
	RCPPS     X3, X3
	SUBPS     X0, X6
	MULPS     X3, X6
	SUBPS     X0, X9
	MULPS     X3, X9
	MOVAPS    X6, X0
	MINPS     X9, X0
	MAXPS     X6, X9
	MINPS     X9, X12
	MOVAPS    X0, X3
	CMPPS     X12, X3, 2 

	RCPPS     X4, X4
	SUBPS     X1, X7
	MULPS     X4, X7
	SUBPS     X1, X10
	MULPS     X4, X10
	MOVAPS    X7, X1
	MINPS     X10, X7
	MAXPS     X1, X10
	MAXPS     X7, X0
	MINPS     X10, X12
	MOVAPS    X0, X4
	CMPPS     X12, X4, 2 

	RCPPS     X5, X5
	SUBPS     X2, X8
	MULPS     X5, X8
	SUBPS     X2, X11
	MULPS     X5, X11
	MOVAPS    X8, X2
	MINPS     X11, X8
	MAXPS     X2, X11
	MAXPS     X8, X0
	MINPS     X11, X12
	MOVAPS    X0, X5
	CMPPS     X12, X5, 2 

	// then we apply the masks to t0 values
	ANDPS     X3, X4
	ANDPS     X4, X5
	ANDPS     X5, X0
	MOVUPS    X0, ret+192(FP)
	RET
