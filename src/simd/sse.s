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

TEXT ·Dot(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	DPPS      $0xf1, X1, X0
	EXTRACTPS $0, X0, AX
	MOVQ      AX, ret+32(FP)
	RET

TEXT ·Cross(SB), NOSPLIT, $0-32
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X1
	MOVAPS    X0, X2
	MOVAPS    X1, X3

	SHUFPS    $0xc9, X0, X0
	SHUFPS    $0xd2, X1, X1
	MULPS     X1, X0

	SHUFPS    $0xd2, X2, X2
	SHUFPS    $0xc9, X3, X3
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

TEXT ·Normalize(SB), NOSPLIT, $0-16
	MOVUPS    a+0(FP), X0 
	MOVAPS    X0, X2
	MULPS     X0, X0
	MOVAPS    X0, X1
	SHUFPS    $147, X0, X0
	ADDPS     X0, X1
	MOVAPS    X1, X0
	SHUFPS    $78, X1, X1
	ADDPS     X1, X0
	RSQRTPS   X0, X0
	MULPS     X2, X0
	MOVUPS    X0, ret+16(FP)
	RET

TEXT ·Normalize4(SB), NOSPLIT, $0-48
	MOVUPS    x4+0(FP), X0
	MOVUPS    y4+16(FP), X1
	MOVUPS    z4+32(FP), X2
	MOVAPS    X0, X3
	MOVAPS    X1, X4
	MOVAPS    X2, X5
	MULPS     X0, X0
	MULPS     X1, X1
	MULPS     X2, X2
	ADDPS     X1, X0
	ADDPS     X2, X0
	RSQRTPS   X0, X0
	MULPS     X0, X3
	MULPS     X0, X4
	MULPS     X0, X5
	MOVUPS    X3, ret+48(FP)
	MOVUPS    X4, ret+64(FP)
	MOVUPS    X5, ret+80(FP)
	RET

// epsilon = 1e-8 in float32 notation
DATA epsilon<>+0x00(SB)/4, $0x322bcc77
DATA epsilon<>+0x04(SB)/4, $0x322bcc77
DATA epsilon<>+0x08(SB)/4, $0x322bcc77
DATA epsilon<>+0x0c(SB)/4, $0x322bcc77
GLOBL epsilon<>(SB), (NOPTR+RODATA), $16

// TODO could be made from epsilon * -0
// only more efficient if -0 calculated on the fly (bitshift?)
DATA minepsilon<>+0x00(SB)/4, $0xb22bcc77
DATA minepsilon<>+0x04(SB)/4, $0xb22bcc77
DATA minepsilon<>+0x08(SB)/4, $0xb22bcc77
DATA minepsilon<>+0x0c(SB)/4, $0xb22bcc77
GLOBL minepsilon<>(SB), (NOPTR+RODATA), $16

DATA one<>+0x00(SB)/4, $0x3f800000
DATA one<>+0x04(SB)/4, $0x3f800000
DATA one<>+0x08(SB)/4, $0x3f800000
DATA one<>+0x0c(SB)/4, $0x3f800000
GLOBL one<>(SB), (NOPTR+RODATA), $16

TEXT ·TriangleIntersect(SB), NOSPLIT, $0-80
	MOVUPS    p0+0(FP), X0
	MOVUPS    p1+16(FP), X1
	MOVUPS    p2+32(FP), X2
	MOVUPS    ro+48(FP), X3
	// p1, p2, ro overwritten 
	SUBPS     X0, X1 // e1
	SUBPS     X0, X2 // e2
	SUBPS     X0, X3 // tvec

	MOVUPS    rd+64(FP), X0
	// rd cross e2
	MOVAPS    X2, X4
	SHUFPS    $0xc9, X0, X0 // 1 2 0 3
	SHUFPS    $0xd2, X4, X4 // 2 0 1 3 
	MULPS     X0, X4
	MOVUPS    rd+64(FP), X0
	MOVAPS    X2, X5
	SHUFPS    $0xd2, X0, X0 // 2 0 1 3
	SHUFPS    $0xc9, X5, X5 // 1 2 0 3
	MULPS     X5, X0
	SUBPS     X0, X4 // pvec

	// e1 dot pvec
	MOVAPS    X1, X5
	DPPS      $0xf1, X4, X5 // det
	MOVAPS    X5, X7
	MOVAPS    minepsilon<>(SB), X0
	CMPPS     X0, X7, 2 // det leq -epsilon
	MOVAPS    epsilon<>(SB), X0
	CMPPS     X5, X0, 2 // epsilon leq det
	ORPS      X0, X7

	RCPPS     X5, X5 // invdet

	// pvec dot tvec * invdet, pvec overwritten
	DPPS      $0xf1, X3, X4
	MULPS     X5, X4 // u
	XORPS     X0, X0
	MOVAPS    X4, X6
	CMPPS     X4, X0, 2 // 0 leq u
	ANDPS     X0, X7
    MOVAPS    one<>(SB), X0
	CMPPS     X0, X6, 2 // u leq 1
	ANDPS     X6, X7

	// tvec cross e1, both overwritten
	MOVAPS    X3, X0
	MOVAPS    X1, X6
	SHUFPS    $0xc9, X3, X3
	SHUFPS    $0xd2, X1, X1
	MULPS     X1, X3
	SHUFPS    $0xd2, X0, X0
	SHUFPS    $0xc9, X6, X6
	MULPS     X6, X0
	SUBPS     X0, X3 // qvec

	MOVUPS    rd+64(FP), X1
	// rd dot qvec * invdet, rd overwritten
	DPPS      $0xf1, X3, X1
	MULPS     X5, X1 // v
	XORPS     X0, X0
	CMPPS     X1, X0, 2 // 0 leq v
	ANDPS     X0, X7
	MOVAPS    one<>(SB), X0
	ADDPS     X4, X1
	CMPPS     X0, X1, 2 // u+v leq 1
	ANDPS     X1, X7

	// e2 dot qvec * invdet
	DPPS      $0xf1, X2, X3
	MULPS     X5, X3 // t
	ANDPS     X7, X3 // apply mask
	MOVD      X3, ret+80(FP)
	RET

TEXT ·Triangle4Intersect(SB), NOSPLIT, $0-240
	MOVUPS    p0x(FP), X0
	MOVUPS    p0y+16(FP), X1
	MOVUPS    p0z+32(FP), X2
	MOVUPS    p1x+48(FP), X3
	MOVUPS    p1y+64(FP), X4
	MOVUPS    p1z+80(FP), X5
	MOVUPS    p2x+96(FP), X6
	MOVUPS    p2y+112(FP), X7
	MOVUPS    p2z+128(FP), X8
	MOVUPS    rox+144(FP), X9
	MOVUPS    roy+160(FP), X10
	MOVUPS    roz+176(FP), X11
	MOVUPS    rdx+192(FP), X12
	MOVUPS    rdy+208(FP), X13
	MOVUPS    rdz+224(FP), X14

	SUBPS     X0, X3
	SUBPS     X1, X4
	SUBPS     X2, X5 // e1
	SUBPS     X0, X6
	SUBPS     X1, X7
	SUBPS     X2, X8 // e2
	SUBPS     X0, X9
	SUBPS     X1, X10
	SUBPS     X2, X11 // tvec

	// rd cross e2
	MOVAPS    X13, X0
	MULPS     X8, X0
	MOVAPS    X14,X15
	MULPS     X7, X15
	SUBPS     X15, X0

	MOVAPS    X14, X1
	MULPS     X6, X1
	MOVAPS    X12, X15
	MULPS     X8, X15
	SUBPS     X15, X1

	MOVAPS    X12, X2
	MULPS     X7, X2
	MOVAPS    X13,X15
	MULPS     X6, X15
	SUBPS     X15, X2 // pvec

	// we can now reuse X12-X15 registers:
	// rd can be read in again when we need it at the end
	// e1 dot pvec
	MOVAPS    X3, X12
	MULPS     X0, X12
	MOVAPS    X4, X15
	MULPS     X1, X15
	ADDPS     X15,X12
	MOVAPS    X5, X15
	MULPS     X2, X15
	ADDPS     X15,X12 // det

	MOVAPS    X12,X13
	MOVAPS    minepsilon<>(SB), X15
	CMPPS     X15,X13, 2 // det leq -epsilon
	MOVAPS    epsilon<>(SB), X15
	CMPPS     X12,X15, 2 // epsilon leq det
	ORPS      X15,X13

	RCPPS     X12,X12 //invdet

	// pvec dot tvec * invdet, pvec overwritten
	MULPS     X9, X0
	MULPS     X10, X1
	MULPS     X11, X2
	ADDPS     X0, X1
	ADDPS     X1, X2
	MULPS     X12, X2 // u
	XORPS     X0, X0
	MOVAPS    X2, X14 // u (since I want X2 for qvecZ)
	CMPPS     X2, X0, 2 // 0 leq u
	ANDPS     X0, X13
	MOVAPS    one<>(SB), X1
	CMPPS     X1, X2, 2 // u leq 1
	ANDPS     X2, X13

	// tvec cross e1, both overwritten
	// tvec in 9-11, e1 in 3-5
	MOVAPS    X10, X0
	MULPS     X5, X0
	MOVAPS    X11, X1
	MULPS     X4, X1
	SUBPS     X1, X0

	MOVAPS    X11, X1
	MULPS     X3, X1
	MOVAPS    X9, X2
	MULPS     X5, X2
	SUBPS     X2, X1

	MOVAPS    X9, X2
	MULPS     X4, X2
	MOVAPS    X10,X15
	MULPS     X3, X15
	SUBPS     X15, X2 // qvec

	// rd dot qvec * invdet
	MOVUPS    rdx+192(FP), X3
	MULPS     X0, X3
	MOVUPS    rdy+208(FP), X4
	MULPS     X1, X4
	ADDPS     X4, X3
	MOVUPS    rdz+224(FP), X4
	MULPS     X2, X4
	ADDPS     X4, X3
	MULPS     X12,X3 // v
	XORPS     X4, X4
	CMPPS     X3, X4, 2 // 0 leq v
	ANDPS     X4, X13
	MOVAPS    one<>(SB), X4
	ADDPS     X14, X3
	CMPPS     X4, X3, 2 // u+v leq 1
	ANDPS     X3, X13

	// e2 dot qvec * invdet
	// e2 in X6-X8, qvec in X0-X2
	MULPS     X6, X0
	MULPS     X7, X1
	MULPS     X8, X2
	ADDPS     X0, X1
	ADDPS     X1, X2
	MULPS     X12, X2 // t
	ANDPS     X13, X2 // apply mask

	MOVUPS    X2, ret+240(FP)
	RET
