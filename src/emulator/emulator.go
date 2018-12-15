package emulator

import (
	"fmt"

	"../instruction"
)

var jumped bool

func signext(ori uint16) uint32 {
	result := uint32(ori)
	result <<= 16
	result = uint32(int32(result) >> 16)
	return result
}

func executeR(instr instruction.RInstruction) bool {
	switch instr.Token {
	case "add":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 + t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "addu":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 + t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "sub":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 - t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "subu":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 - t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "and":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 & t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "or":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 | t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "xor":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 ^ t1
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "nor":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := ^(t0 | t1)
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "slt":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				var t2 uint32 = 0
				if t0 < t1 {
					t2 = 1
				}
				if setRegister(instr.Rd, t2) {
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
					return true
				}
			}
		}
	case "sll":
		if t1, b := getRegister(instr.Rt); b {
			t2 := t1 << instr.Shamt
			if setRegister(instr.Rd, t2) {
				return true
			}
		}
	case "srl":
		if t1, b := getRegister(instr.Rt); b {
			t2 := t1 >> instr.Shamt
			if setRegister(instr.Rd, t2) {
				return true
			}
		}
	case "sra":
		if t1, b := getRegister(instr.Rt); b {
			t2 := uint32(int32(t1) >> instr.Shamt)
			if setRegister(instr.Rd, t2) {
				return true
			}
		}
	case "sllv":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t1 << t0
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "srlv":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t1 >> t0
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "srav":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := uint32(int32(t1) >> t0)
				if setRegister(instr.Rd, t2) {
					return true
				}
			}
		}
	case "jr":
		if t0, b := getRegister(instr.Rs); b {
			current_instr = t0
			jumped = true
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
				return true
			}
		}
	case "addiu":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 + uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				return true
			}
		}
	case "andi":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 & uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				return true
			}
		}
	case "ori":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 | uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				return true
			}
		}
	case "xori":
		if t0, b := getRegister(instr.Rs); b {
			t2 := t0 ^ uint32(instr.Imm)
			if setRegister(instr.Rt, t2) {
				return true
			}
		}
	case "lui":
		t2 := uint32(instr.Imm) << 16
		if setRegister(instr.Rt, t2) {
			return true
		}
	case "lw":
		if t0, b := getRegister(instr.Rs); b {
			if t2, b := memoryRead(t0 + signext(instr.Imm)); b {
				if setRegister(instr.Rt, t2) {
					return true
				}
			}
		}
	case "sw":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				t2 := t0 + signext(instr.Imm)
				if memoryWrite(t2, t1) {
					return true
				}
			}
		}
	case "beq":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				if t0 == t1 {
					current_instr = current_instr + 1 + signext(instr.Imm)
					return true
				}
			}
		}
	case "bne":
		if t0, b := getRegister(instr.Rs); b {
			if t1, b := getRegister(instr.Rt); b {
				if t0 != t1 {
					current_instr = current_instr + 1 + signext(instr.Imm)
					return true
				}
			}
		}
	case "slti":
		if t0, b := getRegister(instr.Rs); b {
			var t2 uint32 = 0
			if t0 < signext(instr.Imm) {
				t2 = 1
			}
			if setRegister(instr.Rt, t2) {
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
				return true
			}
		}
	}
	return false
}

func executeJ(instr instruction.JInstruction) bool {
	switch instr.Token {
	case "j":
		current_instr = instr.Imm
		jumped = true
		return true
	case "jal":
		if setRegister(31, current_instr+1) {
			current_instr = instr.Imm
			jumped = true
			return true
		}
	}
	return false
}

func execute(instr instruction.Instruction) bool {
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

func Execute(entry uint32) bool {
	current_instr = entry
	for {
		jumped = false
		bits, b := memoryRead(current_instr)
		if !b {
			return false
		}
		if bits == 0 {
			break
		}
		instr := instruction.Parse(bits)
		if !execute(instr) {
			return false
		}
		if !jumped {
			current_instr++
		}
	}
	return true
}

func Initialize(bin []uint32) bool {
	jumped = false
	current_instr = 0
	for i := 0; i < 32; i++ {
		setRegister(uint8(i), 0)
	}
	for i := 0; i < MEMORY_SIZE; i++ {
		memoryWrite(uint32(i), 0)
	}
	for i, bits := range bin {
		if !memoryWrite(uint32(i), bits) {
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
