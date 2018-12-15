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

func (this RInstruction) GetToken() string {
	return this.Token
}

func (this RInstruction) ToASM() string {
	return fmt.Sprintf("%-4s $%d, $%d, $%d", this.Token, this.Rd, this.Rs, this.Rt)
}

func (this RInstruction) ToBits() uint32 {
	return (uint32(0x0) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.Rs) & MASK_REG << SHIFT_RS) | (uint32(this.Rt) & MASK_REG << SHIFT_RT) | (uint32(this.Rd) & MASK_REG << SHIFT_RD) | (uint32(this.Shamt) & MASK_SHAMT << SHIFT_SHAMT) | (uint32(this.Funct) & MASK_FUNCT << SHIFT_FUNCT)
}

func CreateR(token string, rs uint8, rt uint8, rd uint8, shamt uint8, funct uint8) RInstruction {
	return RInstruction{token, rs & MASK_REG, rt & MASK_REG, rd & MASK_REG, shamt & MASK_SHAMT, funct & MASK_FUNCT}
}

func ParseR(bits uint32) RInstruction {
	result := CreateR("", uint8(bits>>SHIFT_RS), uint8(bits>>SHIFT_RT), uint8(bits>>SHIFT_RD), uint8(bits>>SHIFT_SHAMT), uint8(bits>>SHIFT_FUNCT))
	if result.Funct == FT_ADD {
		result.Token = "add"
	} else if result.Funct == FT_ADDU {
		result.Token = "addu"
	} else if result.Funct == FT_SUB {
		result.Token = "sub"
	} else if result.Funct == FT_SUBU {
		result.Token = "subu"
	} else if result.Funct == FT_AND {
		result.Token = "and"
	} else if result.Funct == FT_OR {
		result.Token = "or"
	} else if result.Funct == FT_XOR {
		result.Token = "xor"
	} else if result.Funct == FT_NOR {
		result.Token = "nor"
	} else if result.Funct == FT_SLT {
		result.Token = "slt"
	} else if result.Funct == FT_SLTU {
		result.Token = "sltu"
	} else if result.Funct == FT_SLL {
		result.Token = "sll"
	} else if result.Funct == FT_SRL {
		result.Token = "srl"
	} else if result.Funct == FT_SRA {
		result.Token = "sra"
	} else if result.Funct == FT_SLLV {
		result.Token = "sllv"
	} else if result.Funct == FT_SRLV {
		result.Token = "srlv"
	} else if result.Funct == FT_SRAV {
		result.Token = "srav"
	} else if result.Funct == FT_JR {
		result.Token = "jr"
	}
	return result
}

const (
	FT_ADD  = 0x20
	FT_ADDU = 0x21
	FT_SUB  = 0x22
	FT_SUBU = 0x23
	FT_AND  = 0x24
	FT_OR   = 0x25
	FT_XOR  = 0x26
	FT_NOR  = 0x27
	FT_SLT  = 0x2a
	FT_SLTU = 0x2b
	FT_SLL  = 0x00
	FT_SRL  = 0x02
	FT_SRA  = 0x03
	FT_SLLV = 0x04
	FT_SRLV = 0x06
	FT_SRAV = 0x07
	FT_JR   = 0x08
)

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

func Sll(rd uint8, rs uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sll", rs, rt, rd, shamt, FT_SLL)
}

func Srl(rd uint8, rs uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("srl", rs, rt, rd, shamt, FT_SRL)
}

func Sra(rd uint8, rs uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sra", rs, rt, rd, shamt, FT_SRA)
}

func Sllv(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sllv", rs, rt, rd, 0x0, FT_SLLV)
}

func Srlv(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("srlv", rs, rt, rd, 0x0, FT_SRLV)
}

func Srav(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("srav", rs, rt, rd, 0x0, FT_SRAV)
}

func Jr(rs uint8) RInstruction {
	return CreateR("jr", rs, 0x0, 0x0, 0x0, FT_JR)
}
