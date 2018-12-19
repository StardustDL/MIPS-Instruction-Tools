package instruction

import "fmt"

type RInstruction struct {
	Token  string
	Opcode uint8
	Rs     uint8
	Rt     uint8
	Rd     uint8
	Shamt  uint8
	Funct  uint8
}

const (
	FT_MULT    = 0x18
	FT_MULTU   = 0x19
	FT_DIV     = 0x1a
	FT_DIVU    = 0x1b
	FT_MFHI    = 0x10
	FT_MTHI    = 0x11
	FT_MFLO    = 0x12
	FT_MTLO    = 0x13
	FT_ADD     = 0x20
	FT_ADDU    = 0x21
	FT_SUB     = 0x22
	FT_SUBU    = 0x23
	FT_AND     = 0x24
	FT_OR      = 0x25
	FT_XOR     = 0x26
	FT_NOR     = 0x27
	FT_SLT     = 0x2a
	FT_SLTU    = 0x2b
	FT_SLL     = 0x00
	FT_SRL     = 0x02
	FT_SRA     = 0x03
	FT_SLLV    = 0x04
	FT_SRLV    = 0x06
	FT_SRAV    = 0x07
	FT_JR      = 0x08
	FT_JALR    = 0x09
	FT_SYSCALL = 0x0c
	FT_BREAK   = 0x0d
	FT_MUL     = 0x02
)

func (this RInstruction) GetToken() string {
	return this.Token
}

func (this RInstruction) ToASM() string {
	if this.Opcode == OP_SPECIAL {
		if this.Funct == FT_SYSCALL {
			return fmt.Sprintf("%-7s", this.Token)
		} else if this.Funct == FT_BREAK {
			code := uint32(this.Rs)<<15 | uint32(this.Rt)<<10 | uint32(this.Rd)<<5 | uint32(this.Shamt)
			return fmt.Sprintf("%-7s %d", this.Token, code)
		} else if this.Funct == FT_JR {
			return fmt.Sprintf("%-7s $%d", this.Token, this.Rs)
		} else if this.Funct == FT_JALR {
			return fmt.Sprintf("%-7s $%d, $%d", this.Token, this.Rd, this.Rs)
		} else if this.Funct == FT_MULT || this.Funct == FT_MULTU || this.Funct == FT_DIV || this.Funct == FT_DIVU {
			return fmt.Sprintf("%-7s $%d, $%d", this.Token, this.Rs, this.Rt)
		}
	} else if this.Opcode == OP_SPECIAL2 {
		if this.Funct == FT_MUL {
			return fmt.Sprintf("%-7s $%d, $%d, $%d", this.Token, this.Rd, this.Rs, this.Rt)
		}
	}
	return fmt.Sprintf("%-7s $%d, $%d, $%d", this.Token, this.Rd, this.Rs, this.Rt)
}

func (this RInstruction) ToBits() uint32 {
	return (uint32(this.Opcode) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.Rs) & MASK_REG << SHIFT_RS) | (uint32(this.Rt) & MASK_REG << SHIFT_RT) | (uint32(this.Rd) & MASK_REG << SHIFT_RD) | (uint32(this.Shamt) & MASK_SHAMT << SHIFT_SHAMT) | (uint32(this.Funct) & MASK_FUNCT << SHIFT_FUNCT)
}

func CreateR(token string, opcode uint8, rs uint8, rt uint8, rd uint8, shamt uint8, funct uint8) RInstruction {
	return RInstruction{token, opcode & MASK_OPCODE, rs & MASK_REG, rt & MASK_REG, rd & MASK_REG, shamt & MASK_SHAMT, funct & MASK_FUNCT}
}

func ParseR(bits uint32) RInstruction {
	result := CreateR("", uint8(bits>>SHIFT_OPCODE), uint8(bits>>SHIFT_RS), uint8(bits>>SHIFT_RT), uint8(bits>>SHIFT_RD), uint8(bits>>SHIFT_SHAMT), uint8(bits>>SHIFT_FUNCT))
	switch result.Opcode {
	case OP_SPECIAL:
		switch result.Funct {
		case FT_ADD:
			result.Token = "add"
		case FT_ADDU:
			result.Token = "addu"
		case FT_SUB:
			result.Token = "sub"
		case FT_SUBU:
			result.Token = "subu"
		case FT_AND:
			result.Token = "and"
		case FT_OR:
			result.Token = "or"
		case FT_XOR:
			result.Token = "xor"
		case FT_NOR:
			result.Token = "nor"
		case FT_SLT:
			result.Token = "slt"
		case FT_MFHI:
			result.Token = "mfhi"
		case FT_MTHI:
			result.Token = "mthi"
		case FT_MFLO:
			result.Token = "mflo"
		case FT_MTLO:
			result.Token = "mtlo"
		case FT_SLTU:
			result.Token = "sltu"
		case FT_SLL:
			result.Token = "sll"
		case FT_SRL:
			result.Token = "srl"
		case FT_SRA:
			result.Token = "sra"
		case FT_SLLV:
			result.Token = "sllv"
		case FT_SRLV:
			result.Token = "srlv"
		case FT_SRAV:
			result.Token = "srav"
		case FT_JR:
			result.Token = "jr"
		case FT_JALR:
			result.Token = "jalr"
		case FT_SYSCALL:
			result.Token = "syscall"
		case FT_BREAK:
			result.Token = "break"
		case FT_MULT:
			result.Token = "mult"
		case FT_MULTU:
			result.Token = "multu"
		case FT_DIV:
			result.Token = "div"
		case FT_DIVU:
			result.Token = "divu"
		}
	case OP_SPECIAL2:
		switch result.Funct {
		case FT_MUL:
			result.Token = "mul"
		}
	}
	return result
}

func Add(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("add", OP_SPECIAL, rs, rt, rd, 0x0, FT_ADD)
}

func Addu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("addu", OP_SPECIAL, rs, rt, rd, 0x0, FT_ADDU)
}

func Sub(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sub", OP_SPECIAL, rs, rt, rd, 0x0, FT_SUB)
}

func Subu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("subu", OP_SPECIAL, rs, rt, rd, 0x0, FT_SUBU)
}

func Mul(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("mul", OP_SPECIAL2, rs, rt, rd, 0x0, FT_MUL)
}

func And(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("and", OP_SPECIAL, rs, rt, rd, 0x0, FT_AND)
}

func Or(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("or", OP_SPECIAL, rs, rt, rd, 0x0, FT_OR)
}

func Xor(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("xor", OP_SPECIAL, rs, rt, rd, 0x0, FT_XOR)
}

func Nor(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("nor", OP_SPECIAL, rs, rt, rd, 0x0, FT_NOR)
}

func Slt(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("slt", OP_SPECIAL, rs, rt, rd, 0x0, FT_SLT)
}

func Sltu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sltu", OP_SPECIAL, rs, rt, rd, 0x0, FT_SLTU)
}

func Sll(rd uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sll", OP_SPECIAL, 0x00, rt, rd, shamt, FT_SLL)
}

func Srl(rd uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("srl", OP_SPECIAL, 0x00, rt, rd, shamt, FT_SRL)
}

func Sra(rd uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sra", OP_SPECIAL, 0x00, rt, rd, shamt, FT_SRA)
}

func Sllv(rd uint8, rt uint8, rs uint8) RInstruction {
	return CreateR("sllv", OP_SPECIAL, rs, rt, rd, 0x0, FT_SLLV)
}

func Srlv(rd uint8, rt uint8, rs uint8) RInstruction {
	return CreateR("srlv", OP_SPECIAL, rs, rt, rd, 0x0, FT_SRLV)
}

func Srav(rd uint8, rt uint8, rs uint8) RInstruction {
	return CreateR("srav", OP_SPECIAL, rs, rt, rd, 0x0, FT_SRAV)
}

func Jr(rs uint8) RInstruction {
	return CreateR("jr", OP_SPECIAL, rs, 0x0, 0x0, 0x0, FT_JR)
}

func Jalr(rd uint8, rs uint8) RInstruction {
	return CreateR("jalr", OP_SPECIAL, rs, 0x0, rd, 0x0, FT_JALR)
}

func Syscall() RInstruction {
	return CreateR("syscall", OP_SPECIAL, 0x0, 0x0, 0x0, 0x0, FT_SYSCALL)
}

func Break(code uint32) RInstruction {
	code &= 0xfffff
	return CreateR("break", OP_SPECIAL, uint8(code>>15), uint8(code>>10), uint8(code>>5), uint8(code), FT_BREAK)
}

func Mult(rs uint8, rt uint8) RInstruction {
	return CreateR("mult", OP_SPECIAL, rs, rt, 0x0, 0x0, FT_MULT)
}

func Multu(rs uint8, rt uint8) RInstruction {
	return CreateR("multu", OP_SPECIAL, rs, rt, 0x0, 0x0, FT_MULTU)
}

func Div(rs uint8, rt uint8) RInstruction {
	return CreateR("div", OP_SPECIAL, rs, rt, 0x0, 0x0, FT_DIV)
}

func Divu(rs uint8, rt uint8) RInstruction {
	return CreateR("divu", OP_SPECIAL, rs, rt, 0x0, 0x0, FT_DIVU)
}

func Mfhi(rd uint8) RInstruction {
	return CreateR("mfhi", OP_SPECIAL, 0x0, 0x0, rd, 0x0, FT_MFHI)
}

func Mthi(rs uint8) RInstruction {
	return CreateR("mthi", OP_SPECIAL, rs, 0x0, 0x0, 0x0, FT_MFHI)
}

func Mflo(rd uint8) RInstruction {
	return CreateR("mflo", OP_SPECIAL, 0x0, 0x0, rd, 0x0, FT_MFLO)
}

func Mtlo(rs uint8) RInstruction {
	return CreateR("mtlo", OP_SPECIAL, rs, 0x0, 0x0, 0x0, FT_MFLO)
}
