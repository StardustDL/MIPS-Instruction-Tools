package instruction

import "fmt"

type IInstruction struct {
	token  string
	opcode uint8
	rs     uint8
	rt     uint8
	imm    uint16
}

const (
	OP_ADDI  = 0x08
	OP_ADDIU = 0x09
	OP_ANDI  = 0x0c
	OP_ORI   = 0x0d
	OP_XORI  = 0x0e
	OP_LUI   = 0x0f
	OP_LW    = 0x23
	OP_SW    = 0x2b
	OP_BEQ   = 0x04
	OP_BNE   = 0x05
	OP_SLTI  = 0x0a
	OP_SLTIU = 0x0b
)

func (this IInstruction) GetToken() string {
	return this.token
}

func (this IInstruction) ToASM() string {
	if this.opcode == OP_ADDI || this.opcode == OP_ADDIU || this.opcode == OP_ANDI || this.opcode == OP_ORI || this.opcode == OP_XORI || this.opcode == OP_BEQ || this.opcode == OP_BNE || this.opcode == OP_SLTI || this.opcode == OP_SLTIU {
		return fmt.Sprintf("%-4s $%d, $%d, %d", this.token, this.rt, this.rs, this.imm)
	} else if this.opcode == OP_LUI {
		return fmt.Sprintf("%-4s $%d, %d", this.token, this.rt, this.imm)
	} else if this.opcode == OP_LW || this.opcode == OP_SW {
		return fmt.Sprintf("%-4s $%d, %d($%d)", this.token, this.rt, this.imm, this.rs)
	} else {
		return "No this instr"
	}
}

func (this IInstruction) ToBits() uint32 {
	return (uint32(this.opcode) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.rs) & MASK_REG << SHIFT_RS) | (uint32(this.rt) & MASK_REG << SHIFT_RT) | (uint32(this.imm) & MASK_IMM16 << SHIFT_IMMEDIATE)
}

func CreateI(token string, opcode uint8, rs uint8, rt uint8, imm uint16) IInstruction {
	return IInstruction{token, opcode & MASK_OPCODE, rs & MASK_REG, rt & MASK_REG, imm & MASK_IMM16}
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
