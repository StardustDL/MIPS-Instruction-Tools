package emulator

import (
	"fmt"

	"../instruction"
)

func signext(ori uint16) uint32 {
	result := uint32(ori)
	result <<= 16
	result = uint32(int32(result) >> 16)
	return result
}

func advancePC(offset uint32) {
	pc = npc
	npc += offset
}

func executeR(instr instruction.RInstruction) bool {
	switch instr.Token {
	case "add":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 + t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "addu":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 + t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "sub":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 - t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "subu":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 - t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "and":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 & t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "or":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 | t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "xor":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 ^ t1
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "nor":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := ^(t0 | t1)
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "slt":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				var t2 uint32 = 0
				if int32(t0) < int32(t1) {
					t2 = 1
				}
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "sltu":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				var t2 uint32 = 0
				if t0 < t1 {
					t2 = 1
				}
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "sll":
		if t1, b := getRegister(instr.Rt); b {
			t2 := t1 << instr.Shamt
			if setRegister(instr.Rd, t2) {
				advancePC(4)
				return true
			}
		}
	case "srl":
		if t1, b := getRegister(instr.Rt); b {
			t2 := t1 >> instr.Shamt
			if setRegister(instr.Rd, t2) {
				advancePC(4)
				return true
			}
		}
	case "sra":
		if t1, b := getRegister(instr.Rt); b {
			t2 := uint32(int32(t1) >> instr.Shamt)
			if setRegister(instr.Rd, t2) {
				advancePC(4)
				return true
			}
		}
	case "sllv":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t1 << t0
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "srlv":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t1 >> t0
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "srav":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := uint32(int32(t1) >> t0)
				if setRegister(instr.Rd, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "jr":
		if t0, b := getRegister(instr.Rs); b {
			pc = npc
			npc = t0
			return true
		}
	}
	return false
}

func executeI(instr instruction.IInstruction) bool {
	switch instr.Token {
	case "addi":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 + signext(instr.Imm)
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	case "addiu":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 + signext(instr.Imm)
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	case "andi":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 & uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	case "ori":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 | uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	case "xori":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 ^ uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	case "lui":
		t2 := uint32(instr.Imm) << 16
		if setRegister(instr.Rt, t2) {
			advancePC(4)
			return true
		}
	case "lw":
		if t0, b := getRegister(instr.Rs); b {
			if t2, b := memoryRead(t0+signext(instr.Imm), 4); b {
				if setRegister(instr.Rt, t2) {
					advancePC(4)
					return true
				}
			}
		}
	case "sw":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 + signext(instr.Imm)
				if memoryWrite(t2, 4, t1) {
					advancePC(4)
					return true
				}
			}
		}
	case "beq":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				if t0 == t1 {
					advancePC(signext(instr.Imm) << 2)
				} else {
					advancePC(4)
				}
				return true
			}
		}
	case "bne":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				if t0 != t1 {
					advancePC(signext(instr.Imm) << 2)
				} else {
					advancePC(4)
				}
				return true
			}
		}
	case "bgez":
		if t0, b := getRegister(instr.Rs); b {
			if t0 >= 0 {
				advancePC(signext(instr.Imm) << 2)
			} else {
				advancePC(4)
			}
			return true
		}
	case "bgezal":
		if t0, b := getRegister(instr.Rs); b {
			if t0 >= 0 {
				if !setRegister(31, npc+4) {
					return false
				}
				advancePC(signext(instr.Imm) << 2)
			} else {
				advancePC(4)
			}
			return true
		}
	case "bgtz":
		if t0, b := getRegister(instr.Rs); b {
			if t0 > 0 {
				advancePC(signext(instr.Imm) << 2)
			} else {
				advancePC(4)
			}
			return true
		}
	case "bltz":
		if t0, b := getRegister(instr.Rs); b {
			if t0 < 0 {
				advancePC(signext(instr.Imm) << 2)
			} else {
				advancePC(4)
			}
			return true
		}
	case "bltzal":
		if t0, b := getRegister(instr.Rs); b {
			if t0 < 0 {
				if !setRegister(31, npc+4) {
					return false
				}
				advancePC(signext(instr.Imm) << 2)
			} else {
				advancePC(4)
			}
			return true
		}
	case "blez":
		if t0, b := getRegister(instr.Rs); b {
			if t0 <= 0 {
				advancePC(signext(instr.Imm) << 2)
			} else {
				advancePC(4)
			}
			return true
		}
	case "slti":
		if t0, b := getRegister(instr.Rs); b {
			var t2 uint32 = 0
			if int32(t0) < int32(signext(instr.Imm)) {
				t2 = 1
			}
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	case "sltiu":
		if t0, b := getRegister(instr.Rs); b {
			var t2 uint32 = 0
			if t0 < uint32(instr.Imm) {
				t2 = 1
			}
			if setRegister(instr.Rt, t2) {
				advancePC(4)
				return true
			}
		}
	}
	return false
}

func executeJ(instr instruction.JInstruction) bool {
	switch instr.Token {
	case "j":
		pc = npc
		npc = (pc & 0xf0000000) | (instr.Imm << 2)
		return true
	case "jal":
		if setRegister(31, npc+4) {
			pc = npc
			npc = (pc & 0xf0000000) | (instr.Imm << 2)
			return true
		}
	}
	return false
}

func execute(instr instruction.Instruction) bool {
	fmt.Printf("%x: %08x %s\n", pc, instr.ToBits(), instr.ToASM())
	switch instr.(type) {
	case instruction.RInstruction:
		return executeR(instr.(instruction.RInstruction))
	case instruction.IInstruction:
		return executeI(instr.(instruction.IInstruction))
	case instruction.JInstruction:
		return executeJ(instr.(instruction.JInstruction))
	}
	return false
}

const END_INSTR uint32 = 0xffffffff

func Execute(entry uint32) bool {
	pc = entry
	npc = pc + 4
	for {
		bits, b := memoryRead(pc, 4)
		if !b {
			return false
		}
		if bits == END_INSTR {
			break
		}
		instr := instruction.Parse(bits)
		if !execute(instr) {
			return false
		}
	}
	return true
}

func Initialize(bin []uint8) bool {
	pc = 0
	npc = 0
	for i := 0; i < 32; i++ {
		setRegister(uint8(i), 0)
	}
	for i := uint32(0); i < MEMORY_SIZE; i++ {
		memoryWrite(uint32(i), 1, 0)
	}
	for i, bits := range bin {
		if !memoryWrite(uint32(i), 1, uint32(bits)) {
			return false
		}
	}
	return true
}

func ShowRegisters() {
	for i := 0; i < 32; i++ {
		val, _ := getRegister(uint8(i))
		fmt.Printf("$%d %x\n", i, val)
	}
}
