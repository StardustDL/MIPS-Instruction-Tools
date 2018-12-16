package instruction

import "fmt"

type JInstruction struct {
	Token  string
	Opcode uint8
	Imm    uint32
}

const (
	OP_J   = 0x02
	OP_JAL = 0x03
)

func (this JInstruction) GetToken() string {
	return this.Token
}

func (this JInstruction) ToASM() string {
	return fmt.Sprintf("%-5s 0x%x", this.Token, this.Imm<<2)
}

func (this JInstruction) ToBits() uint32 {
	return (uint32(this.Opcode) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.Imm) & MASK_IMM26 << SHIFT_IMMEDIATE)
}

func CreateJ(token string, opcode uint8, imm uint32) JInstruction {
	return JInstruction{token, opcode & MASK_OPCODE, (imm & MASK_IMM26)}
}

func ParseJ(bits uint32) JInstruction {
	result := CreateJ("", uint8(bits>>SHIFT_OPCODE), uint32(bits>>SHIFT_IMMEDIATE))
	if result.Opcode == OP_J {
		result.Token = "j"
	} else if result.Opcode == OP_JAL {
		result.Token = "jal"
	}
	return result
}

func J(imm uint32) JInstruction {
	return CreateJ("j", OP_J, imm>>2)
}

func Jal(imm uint32) JInstruction {
	return CreateJ("jal", OP_JAL, imm>>2)
}
