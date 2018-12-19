package instruction

import (
	"errors"
	"fmt"
)

type IInstruction struct {
	Token  string
	Opcode uint8
	Rs     uint8
	Rt     uint8
	Imm    uint16
}

const (
	OP_ADDI     = 0x08
	OP_ADDIU    = 0x09
	OP_ANDI     = 0x0c
	OP_ORI      = 0x0d
	OP_XORI     = 0x0e
	OP_LUI      = 0x0f
	OP_LB       = 0x20
	OP_LH       = 0x21
	OP_LW       = 0x23
	OP_LBU      = 0x24
	OP_LHU      = 0x25
	OP_SB       = 0x28
	OP_SH       = 0x29
	OP_SW       = 0x2b
	OP_BEQ      = 0x04
	OP_BNE      = 0x05
	OP_BGEZ     = 0x01
	OP_BGEZAL   = 0x01
	OP_BLTZ     = 0x01
	OP_BLTZAL   = 0x01
	OP_BGTZ     = 0x07
	OP_BLEZ     = 0x06
	OP_SLTI     = 0x0a
	OP_SLTIU    = 0x0b
	OP_SPECIAL  = 0x00
	OP_SPECIAL2 = 0x1c
)

func (this IInstruction) GetToken() string {
	return this.Token
}

func (this IInstruction) ToASM() string {
	if this.Opcode == OP_ADDI || this.Opcode == OP_ADDIU || this.Opcode == OP_ANDI || this.Opcode == OP_ORI || this.Opcode == OP_XORI || this.Opcode == OP_SLTI || this.Opcode == OP_SLTIU {
		return fmt.Sprintf("%-7s $%d, $%d, 0x%x", this.Token, this.Rt, this.Rs, this.Imm)
	} else if this.Opcode == OP_BEQ || this.Opcode == OP_BNE {
		return fmt.Sprintf("%-7s $%d, $%d, 0x%x", this.Token, this.Rs, this.Rt, this.Imm)
	} else if this.Opcode == OP_LUI {
		return fmt.Sprintf("%-7s $%d, 0x%x", this.Token, this.Rt, this.Imm)
	} else if this.Opcode == OP_LW || this.Opcode == OP_SW || this.Opcode == OP_LB || this.Opcode == OP_SB || this.Opcode == OP_LBU || this.Opcode == OP_LH || this.Opcode == OP_LHU || this.Opcode == OP_SH {
		return fmt.Sprintf("%-7s $%d, 0x%x($%d)", this.Token, this.Rt, this.Imm, this.Rs)
	} else if this.Opcode == OP_BGEZ || this.Opcode == OP_BGEZAL || this.Opcode == OP_BLTZ || this.Opcode == OP_BLTZAL || this.Opcode == OP_BLEZ || this.Opcode == OP_BGTZ {
		return fmt.Sprintf("%-7s $%d, 0x%x", this.Token, this.Rs, this.Imm)
	} else {
		panic(errors.New(fmt.Sprintf("No this instr %d", this.Opcode)))
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
	case OP_LB:
		result.Token = "lb"
	case OP_LBU:
		result.Token = "lbu"
	case OP_LH:
		result.Token = "lh"
	case OP_LHU:
		result.Token = "lhu"
	case OP_SB:
		result.Token = "sb"
	case OP_LW:
		result.Token = "lw"
	case OP_SW:
		result.Token = "sw"
	case OP_SH:
		result.Token = "sh"
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

func Sh(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("sh", OP_SH, rs, rt, imm)
}

func Lb(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("lb", OP_LB, rs, rt, imm)
}

func Lbu(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("lbu", OP_LBU, rs, rt, imm)
}

func Lh(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("lh", OP_LH, rs, rt, imm)
}

func Lhu(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("lhu", OP_LHU, rs, rt, imm)
}

func Sb(rt uint8, rs uint8, imm uint16) IInstruction {
	return CreateI("sb", OP_SB, rs, rt, imm)
}

func Beq(rs uint8, rt uint8, imm uint16) IInstruction {
	return CreateI("beq", OP_BEQ, rs, rt, imm)
}

func Bne(rs uint8, rt uint8, imm uint16) IInstruction {
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
