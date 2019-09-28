TEXT ·Add(SB),$0
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X2
	ADDPS     X2, X0
	MOVUPS    X0, ret+32(FP)
	RET

TEXT ·Sub(SB),$0
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X2
	SUBPS     X2, X0
	MOVUPS    X0, ret+32(FP)
	RET

TEXT ·Min(SB),$0
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X2
	MINPS     X2, X0
	MOVUPS    X0, ret+32(FP)
	RET

TEXT ·Max(SB),$0
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X2
	MAXPS     X2, X0
	MOVUPS    X0, ret+32(FP)
	RET
	
TEXT ·Dot(SB),$0
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X2
	MULPS     X2, X0
	HADDPS    X0, X0
	HADDPS    X0, X0
	MOVUPS    X0, ret+32(FP)
	RET
