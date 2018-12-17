package main

import (
	"fmt"

	emu "../emulator"
	ins "../instruction"
)

var instrs []ins.Instruction

const (
	PORT_INPUT  = 0x0000
	PORT_OUTPUT = 0x0080
	PORT_LED    = 0x0100
	PORT_DIGIT  = 0x0101
	SEG_TEXT    = 0x1000
	ENTRY       = 0x1000
)

var memory [emu.MEMORY_SIZE]uint8

func initMemory(text []uint32) {
	for i := uint32(0); i < emu.MEMORY_SIZE; i++ {
		memory[i] = 0
	}

	for i, bits := range text {
		for j := uint32(0); j < 4; j++ {
			memory[SEG_TEXT+(uint32(i)<<2)+j] = uint8(bits >> (j << 3) & 0xff)
		}
	}
	emu.Initialize(memory[:])
}

func createSymbolInText() uint32 {
	return SEG_TEXT + uint32(len(instrs))
}

func ai(instr ins.Instruction) {
	instrs = append(instrs, instr)
}

func main() {

	fmt.Println("Running")

	fmt.Println("Emulating")
	initMemory(ins.ToBin(instrs))
	emu.Execute(ENTRY)
}
