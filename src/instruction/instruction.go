package instruction

type Instruction interface {
	GetToken() string
	ToASM() string
	ToBits() uint32
}

const (
	SHIFT_FUNCT     = 0
	SHIFT_SHAMT     = 6
	SHIFT_RD        = 11
	SHIFT_RT        = 16
	SHIFT_RS        = 21
	SHIFT_OPCODE    = 26
	SHIFT_IMMEDIATE = 0
)

const (
	MASK_REG    = 0x1f
	MASK_SHAMT  = 0x1f
	MASK_FUNCT  = 0x3f
	MASK_OPCODE = 0x3f
	MASK_IMM16  = 0xffff
	MASK_IMM26  = 0x3ffffff
)

const (
	GPR_ZERO = uint8(iota)
	GPR_AT
	GPR_V0
	GPR_V1
	GPR_A0
	GPR_A1
	GPR_A2
	GPR_A3
	GPR_T0
	GPR_T1
	GPR_T2
	GPR_T3
	GPR_T4
	GPR_T5
	GPR_T6
	GPR_T7
	GPR_S0
	GPR_S1
	GPR_S2
	GPR_S3
	GPR_S4
	GPR_S5
	GPR_S6
	GPR_S7
	GPR_T8
	GPR_T9
	GPR_K0
	GPR_K1
	GPR_GP
	GPR_SP
	GPR_FP
	GPR_RA
)

func Parse(bits uint32) Instruction {
	opcode := bits >> SHIFT_OPCODE
	if opcode == OP_SPECIAL || opcode == OP_SPECIAL2 {
		return ParseR(bits)
	} else {
		var result Instruction = ParseI(bits)
		if result.GetToken() != "" {
			return result
		}

		result = ParseJ(bits)

		return result
	}
}

func ToBin(instrs []Instruction) []uint32 {
	result := make([]uint32, 0, len(instrs))
	for _, instr := range instrs {
		result = append(result, instr.ToBits())
	}
	return result
}
