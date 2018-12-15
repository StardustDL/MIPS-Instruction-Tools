package main

import (
	"fmt"
	"strconv"

	"./instruction"
)

func createTestInstructions() []instruction.Instruction {
	instrs := make([]instruction.Instruction, 0, 1024)
	instrs = append(instrs, instruction.Lw(2, 0, 0x0))
	instrs = append(instrs, instruction.Lw(3, 0, 0x1))
	instrs = append(instrs, instruction.Add(1, 2, 3))
	return instrs
}

func outputASMs(instrs []instruction.Instruction) {
	for _, instr := range instrs {
		fmt.Println(instr.ToASM())
	}
}

func outputBits(instrs []instruction.Instruction) {
	for _, instr := range instrs {
		fmt.Printf("%032s\n", strconv.FormatUint(uint64(instr.ToBits()), 2))
	}
}

func main() {
	instrs := createTestInstructions()
	fmt.Println("Asm:")
	outputASMs(instrs)
	fmt.Println()
	fmt.Println("Bit:")
	outputBits(instrs)
}
