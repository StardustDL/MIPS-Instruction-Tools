package instruction

import "fmt"

type RInstruction struct {
	Token string
	Rs    uint8
	Rt    uint8
	Rd    uint8
	Shamt uint8
	Funct uint8
}

const (
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
	FT_SYSCALL = 0x0a
)

func (this RInstruction) GetToken() string {
	return this.Token
}

func (this RInstruction) ToASM() string {
	if this.Token == "nop" || this.Token == "syscall" {
		return fmt.Sprintf("%-7s", this.Token)
	} else if this.Token == "jr" {
		return fmt.Sprintf("%-7s $%d", this.Token, this.Rs)
	}
	return fmt.Sprintf("%-7s $%d, $%d, $%d", this.Token, this.Rd, this.Rs, this.Rt)
}

func (this RInstruction) ToBits() uint32 {
	return (uint32(0x0) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.Rs) & MASK_REG << SHIFT_RS) | (uint32(this.Rt) & MASK_REG << SHIFT_RT) | (uint32(this.Rd) & MASK_REG << SHIFT_RD) | (uint32(this.Shamt) & MASK_SHAMT << SHIFT_SHAMT) | (uint32(this.Funct) & MASK_FUNCT << SHIFT_FUNCT)
}

func CreateR(token string, rs uint8, rt uint8, rd uint8, shamt uint8, funct uint8) RInstruction {
	return RInstruction{token, rs & MASK_REG, rt & MASK_REG, rd & MASK_REG, shamt & MASK_SHAMT, funct & MASK_FUNCT}
}

func ParseR(bits uint32) RInstruction {
	if bits == 0 {
		return Nop()
	}
	result := CreateR("", uint8(bits>>SHIFT_RS), uint8(bits>>SHIFT_RT), uint8(bits>>SHIFT_RD), uint8(bits>>SHIFT_SHAMT), uint8(bits>>SHIFT_FUNCT))
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
	case FT_SYSCALL:
		result.Token = "syscall"
	}
	return result
}

func Add(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("add", rs, rt, rd, 0x0, FT_ADD)
}

func Addu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("addu", rs, rt, rd, 0x0, FT_ADDU)
}

func Sub(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sub", rs, rt, rd, 0x0, FT_SUB)
}

func Subu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("subu", rs, rt, rd, 0x0, FT_SUBU)
}

func And(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("and", rs, rt, rd, 0x0, FT_AND)
}

func Or(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("or", rs, rt, rd, 0x0, FT_OR)
}

func Xor(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("xor", rs, rt, rd, 0x0, FT_XOR)
}

func Nor(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("nor", rs, rt, rd, 0x0, FT_NOR)
}

func Slt(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("slt", rs, rt, rd, 0x0, FT_SLT)
}

func Sltu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sltu", rs, rt, rd, 0x0, FT_SLTU)
}

func Sll(rd uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sll", 0x00, rt, rd, shamt, FT_SLL)
}

func Srl(rd uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("srl", 0x00, rt, rd, shamt, FT_SRL)
}

func Sra(rd uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sra", 0x00, rt, rd, shamt, FT_SRA)
}

func Sllv(rd uint8, rt uint8, rs uint8) RInstruction {
	return CreateR("sllv", rs, rt, rd, 0x0, FT_SLLV)
}

func Srlv(rd uint8, rt uint8, rs uint8) RInstruction {
	return CreateR("srlv", rs, rt, rd, 0x0, FT_SRLV)
}

func Srav(rd uint8, rt uint8, rs uint8) RInstruction {
	return CreateR("srav", rs, rt, rd, 0x0, FT_SRAV)
}

func Jr(rs uint8) RInstruction {
	return CreateR("jr", rs, 0x0, 0x0, 0x0, FT_JR)
}

func Nop() RInstruction {
	return CreateR("nop", 0x0, 0x0, 0x0, 0x0, 0x0)
}

func Syscall() RInstruction {
	return CreateR("syscall", 0x0, 0x0, 0x0, 0x0, FT_SYSCALL)
}
