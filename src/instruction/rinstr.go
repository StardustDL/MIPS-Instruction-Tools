package instruction

import "fmt"

type RInstruction struct {
	token string
	rs    uint8
	rt    uint8
	rd    uint8
	shamt uint8
	funct uint8
}

func (this RInstruction) GetToken() string {
	return this.token
}

func (this RInstruction) ToASM() string {
	return fmt.Sprintf("%-4s $%d, $%d, $%d", this.token, this.rd, this.rs, this.rt)
}

func (this RInstruction) ToBits() uint32 {
	return (uint32(0x0) & MASK_OPCODE << SHIFT_OPCODE) | (uint32(this.rs) & MASK_REG << SHIFT_RS) | (uint32(this.rt) & MASK_REG << SHIFT_RT) | (uint32(this.rd) & MASK_REG << SHIFT_RD) | (uint32(this.shamt) & MASK_SHAMT << SHIFT_SHAMT) | (uint32(this.funct) & MASK_FUNCT << SHIFT_FUNCT)
}

func CreateR(token string, rs uint8, rt uint8, rd uint8, shamt uint8, funct uint8) RInstruction {
	return RInstruction{token, rs & MASK_REG, rt & MASK_REG, rd & MASK_REG, shamt & MASK_SHAMT, funct & MASK_FUNCT}
}

func Add(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("add", rs, rt, rd, 0x0, 0x20)
}

func Addu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("addu", rs, rt, rd, 0x0, 0x21)
}

func Sub(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sub", rs, rt, rd, 0x0, 0x22)
}

func Subu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("subu", rs, rt, rd, 0x0, 0x23)
}

func And(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("and", rs, rt, rd, 0x0, 0x24)
}

func Or(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("or", rs, rt, rd, 0x0, 0x25)
}

func Xor(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("xor", rs, rt, rd, 0x0, 0x26)
}

func Nor(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("nor", rs, rt, rd, 0x0, 0x27)
}

func Slt(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("slt", rs, rt, rd, 0x0, 0x2a)
}

func Sltu(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sltu", rs, rt, rd, 0x0, 0x2b)
}

func Sll(rd uint8, rs uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sll", rs, rt, rd, shamt, 0x00)
}

func Srl(rd uint8, rs uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("srl", rs, rt, rd, shamt, 0x02)
}

func Sra(rd uint8, rs uint8, rt uint8, shamt uint8) RInstruction {
	return CreateR("sra", rs, rt, rd, shamt, 0x03)
}

func Sllv(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("sllv", rs, rt, rd, 0x0, 0x04)
}

func Srlv(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("srlv", rs, rt, rd, 0x0, 0x06)
}

func Srav(rd uint8, rs uint8, rt uint8) RInstruction {
	return CreateR("srav", rs, rt, rd, 0x0, 0x07)
}

func Jr(rs uint8) RInstruction {
	return CreateR("jr", rs, 0x0, 0x0, 0x0, 0x08)
}
