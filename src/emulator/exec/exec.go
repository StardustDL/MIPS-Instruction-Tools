package exec

import (
	"../../instruction"
	"../cpu"
)

type (
	instr  = instruction.Instruction
	rinstr = instruction.RInstruction
	iinstr = instruction.IInstruction
	jinstr = instruction.JInstruction
)

type ExecRFunc func(it rinstr)
type ExecIFunc func(it iinstr)
type ExecJFunc func(it jinstr)
type SignalHandler func(code uint32)

var ExecTable map[string]interface{}

var jumped bool

var npc uint32

var State uint32

var IsDebug bool

var breakHandler, syscallHandler SignalHandler

const (
	MEMU_INITIALIZED = uint32(iota)
	MEMU_RUNNING
	MEMU_EXITED
	MEMU_ERROR
)

func advancePC(offset uint32) {
	cpu.PC = npc
	npc += offset
}

func jumpOneDelay(addr uint32) {
	cpu.PC = npc
	npc = addr
	jumped = true
}

func beforeCall() {
	cpu.SetGPR(instruction.GPR_RA, npc+4)
}

func beforeCallSet(rd uint8) {
	cpu.SetGPR(rd, npc+4)
}

func signext(ori uint32, len uint8) uint32 {
	if !(len == 1 || len == 2 || len == 4) {
		panic("signext: Len error")
	}
	result := uint32(ori)
	result <<= len << 3
	result = uint32(int32(result) >> (len << 3))
	return result
}

func signext8(ori uint8) uint32 {
	return signext(uint32(ori), 1)
}

func signext16(ori uint16) uint32 {
	return signext(uint32(ori), 2)
}

func signext64(ori uint32) uint64 {
	result := uint64(ori)
	result <<= 32
	result = uint64(int64(result) >> 32)
	return result
}

func UpdatePC() {
	if !jumped {
		advancePC(4)
	} else {
		jumped = false
	}
}

func InitializeTable(breakH SignalHandler,syscallH SignalHandler) {
	ExecTable = map[string]interface{}{
		"add":     ExecRFunc(add),
		"addu":    ExecRFunc(addu),
		"addi":    ExecIFunc(addi),
		"addiu":   ExecIFunc(addiu),
		"sub":     ExecRFunc(sub),
		"subu":    ExecRFunc(subu),
		"lui":     ExecIFunc(lui),
		"mul":     ExecRFunc(mul),
		"mult":    ExecRFunc(mult),
		"multu":   ExecRFunc(multu),
		"div":     ExecRFunc(div),
		"divu":    ExecRFunc(divu),
		"mfhi":    ExecRFunc(mfhi),
		"mflo":    ExecRFunc(mflo),
		"mthi":    ExecRFunc(mthi),
		"mtlo":    ExecRFunc(mtlo),
		"slt":     ExecRFunc(slt),
		"slti":    ExecIFunc(slti),
		"sltu":    ExecRFunc(sltu),
		"sltiu":   ExecIFunc(sltiu),
		"sll":     ExecRFunc(sll),
		"sllv":    ExecRFunc(sllv),
		"sra":     ExecRFunc(sra),
		"srav":    ExecRFunc(srav),
		"srl":     ExecRFunc(srl),
		"srlv":    ExecRFunc(srlv),
		"and":     ExecRFunc(and),
		"andi":    ExecIFunc(andi),
		"or":      ExecRFunc(or),
		"ori":     ExecIFunc(ori),
		"nor":     ExecRFunc(nor),
		"xor":     ExecRFunc(xor),
		"xori":    ExecIFunc(xori),
		"lb":      ExecIFunc(lb),
		"lbu":     ExecIFunc(lbu),
		"lh":      ExecIFunc(lh),
		"lhu":     ExecIFunc(lhu),
		"lw":      ExecIFunc(lw),
		"sb":      ExecIFunc(sb),
		"sh":      ExecIFunc(sh),
		"sw":      ExecIFunc(sw),
		"beq":     ExecIFunc(beq),
		"bne":     ExecIFunc(bne),
		"bgez":    ExecIFunc(bgez),
		"bgezal":  ExecIFunc(bgezal),
		"bgtz":    ExecIFunc(bgtz),
		"blez":    ExecIFunc(blez),
		"bltz":    ExecIFunc(bltz),
		"bltzal":  ExecIFunc(bltzal),
		"j":       ExecJFunc(j),
		"jal":     ExecJFunc(jal),
		"jr":      ExecRFunc(jr),
		"jalr":    ExecRFunc(jalr),
		"syscall": ExecRFunc(syscall),
		"break":   ExecRFunc(_break),
	}
	breakHandler = breakH
	syscallHandler = syscallH
}

func InitializePC(pc uint32) {
	jumped = false
	cpu.PC = pc
	npc = pc + 4
}
