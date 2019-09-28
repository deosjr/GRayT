TEXT ·Add(SB),$0
	MOVUPS    a+0(FP), X0
	MOVUPS    b+16(FP), X2
	ADDPS     X2, X0
	MOVUPS    X0, ret+32(FP)
	RET
