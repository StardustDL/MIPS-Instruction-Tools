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

func Parse(bits uint32) Instruction {
	if bits>>SHIFT_OPCODE == 0 {
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
