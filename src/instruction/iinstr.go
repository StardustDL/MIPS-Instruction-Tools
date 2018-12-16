package instruction

import "fmt"

type IInstruction struct {
	Token  string
	Opcode uint8
	Rs     uint8
	Rt     uint8
	Imm    uint16
}

const (
	OP_ADDI   = 0x08
	OP_ADDIU  = 0x09
	OP_ANDI   = 0x0c
	OP_ORI    = 0x0d
	OP_XORI   = 0x0e
	OP_LUI    = 0x0f
	OP_LW     = 0x23
	OP_SW     = 0x2b
	OP_BEQ    = 0x04
	OP_BNE    = 0x05
	OP_BGEZ   = 0x01
	OP_BGEZAL = 0x01
	OP_BLTZ   = 0x01
	OP_BLTZAL = 0x01
	OP_BGTZ   = 0x07
	OP_BLEZ   = 0x06
	OP_SLTI   = 0x0a
	OP_SLTIU  = 0x0b
)

func (this IInstruction) GetToken() string {
	return this.Token
}

func (this IInstruction) ToASM() string {
	if this.Opcode == OP_ADDI || this.Opcode == OP_ADDIU || this.Opcode == OP_ANDI || this.Opcode == OP_ORI || this.Opcode == OP_XORI || this.Opcode == OP_BEQ || this.Opcode == OP_BNE || this.Opcode == OP_SLTI || this.Opcode == OP_SLTIU {
		return fmt.Sprintf("%-4s $%d, $%d, %d", this.Token, this.Rt, this.Rs, this.Imm)
	} else if this.Opcode == OP_LUI {
		return fmt.Sprintf("%-4s $%d, %d", this.Token, this.Rt, this.Imm)
	} else if this.Opcode == OP_LW || this.Opcode == OP_SW {
		return fmt.Sprintf("%-4s $%d, %d($%d)", this.Token, this.Rt, this.Imm, this.Rs)
	} else if this.Opcode == OP_BGEZ || this.Opcode == OP_BGEZAL || this.Opcode == OP_BLTZ || this.Opcode == OP_BLTZAL || this.Opcode == OP_BLEZ || this.Opcode == OP_BGTZ {
		return fmt.Sprintf("%-4s $%d, %d", this.Token, this.Rs, this.Imm)
	} else {
		return "No this instr"
	}
}

func (this IInstruction) ToBits() uint32 {
	return (uint32(this.Opcode) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.Rs) & MASK_REG << SHIFT_RS) | (uint32(this.Rt) & MASK_REG << SHIFT_RT) | (uint32(this.Imm) & MASK_IMM16 << SHIFT_IMMEDIATE)
}

func CreateI(token string, Opcode uint8, rs uint8, rt uint8, imm uint16) IInstruction {
	return IInstruction{token, Opcode & MASK_OPCODE, rs & MASK_REG, rt & MASK_REG, imm & MASK_IMM16}
}

func ParseI(bits uint32) IInstruction {
	result := CreateI("", uint8(bits>>SHIFT_OPCODE), uint8(bits>>SHIFT_RS), uint8(bits>>SHIFT_RT), uint16(bits>>SHIFT_IMMEDIATE))
	switch result.Opcode {
	case OP_ADDI:
		result.Token = "addi"
	case OP_ADDIU:
		result.Token = "addiu"
	case OP_ANDI:
		result.Token = "andi"
	case OP_ORI:
		result.Token = "ori"
	case OP_XORI:
		result.Token = "xori"
	case OP_LUI:
		result.Token = "lui"
	case OP_LW:
		result.Token = "lw"
	case OP_SW:
		result.Token = "sw"
	case OP_BEQ:
		result.Token = "beq"
	case OP_BNE:
		result.Token = "bne"
	case OP_SLTI:
		result.Token = "slti"
	case OP_SLTIU:
		result.Token = "sltiu"
	case OP_BLEZ:
		result.Token = "blez"
	case OP_BGTZ:
		result.Token = "bgtz"
	case OP_BGEZ:
		switch result.Rt {
		case 0x01:
			result.Token = "bgez"
		case 0x11:
			result.Token = "bgezal"
		case 0x00:
			result.Token = "bltz"
		case 0x20:
			result.Token = "bltzal"
		}
	}
	return result
}

func Addi(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("addi", OP_ADDI, rs, rt, imm)
}

func Addiu(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("addiu", OP_ADDIU, rs, rt, imm)
}

func Andi(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("andi", OP_ANDI, rs, rt, imm)
}

func Ori(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("ori", OP_ORI, rs, rt, imm)
}

func Xori(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("xori", OP_XORI, rs, rt, imm)
}

func Lui(rt uint8, imm uint16) IInstruction {
	return CreateI("lui", OP_LUI, 0x0, rt, imm)
}

func Lw(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("lw", OP_LW, rs, rt, imm)
}

func Sw(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("sw", OP_SW, rs, rt, imm)
}

func Beq(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("beq", OP_BEQ, rs, rt, imm)
}

func Bne(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("bne", OP_BNE, rs, rt, imm)
}

func Slti(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("slti", OP_SLTI, rs, rt, imm)
}

func Sltiu(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("sltiu", OP_SLTIU, rs, rt, imm)
}

func Blez(rs uint8, imm uint16) IInstruction {
	return CreateI("blez", OP_BLEZ, rs, 0x00, imm)
}

func Bgtz(rs uint8, imm uint16) IInstruction {
	return CreateI("bgez", OP_BGTZ, rs, 0x00, imm)
}

func Bgez(rs uint8, imm uint16) IInstruction {
	return CreateI("bgez", OP_BGEZ, rs, 0x01, imm)
}

func Bgezal(rs uint8, imm uint16) IInstruction {
	return CreateI("bgezal", OP_BGEZAL, rs, 0x11, imm)
}

func Bltz(rs uint8, imm uint16) IInstruction {
	return CreateI("bltz", OP_BLTZ, rs, 0x00, imm)
}

func Bltzal(rs uint8, imm uint16) IInstruction {
	return CreateI("bltzal", OP_BLTZAL, rs, 0x20, imm)
}
