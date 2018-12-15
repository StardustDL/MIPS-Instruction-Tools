package instruction

import "fmt"

type JInstruction struct {
	token  string
	opcode uint8
	imm    uint32
}

func (this JInstruction) GetToken() string {
	return this.token
}

func (this JInstruction) ToASM() string {
	return fmt.Sprintf("%-4s %d", this.token, this.imm)
}

func (this JInstruction) ToBits() uint32 {
	return (uint32(this.opcode) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.imm) & MASK_IMM26 << SHIFT_IMMEDIATE)
}

func CreateJ(token string, opcode uint8, imm uint32) JInstruction {
	return JInstruction{token, opcode & MASK_OPCODE, imm & MASK_IMM26}
}

func J(imm uint32) JInstruction {
	return CreateJ("j", 0x02, imm)
}

func Jal(imm uint32) JInstruction {
	return CreateJ("jal", 0x03, imm)
}
