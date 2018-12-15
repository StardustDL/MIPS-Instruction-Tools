package main

import (
	"fmt"
	"strconv"

	"./emulator"
	"./instruction"
)

func createTestInstructions() []instruction.Instruction {
	instrs := make([]instruction.Instruction, 0, 1024)
	instrs = append(instrs, instruction.Lw(2, 0, 0x0))
	instrs = append(instrs, instruction.Lw(3, 0, 0x1))
	instrs = append(instrs, instruction.Add(1, 2, 3))
	return instrs
}

func toBin(instrs []instruction.Instruction) []uint32 {
	result := make([]uint32, 0, len(instrs))
	for _, instr := range instrs {
		result = append(result, instr.ToBits())
	}
	return result
}

func testToAndParse(instrs []instruction.Instruction) {
	for _, instr := range instrs {
		to := instr.ToBits()
		parsed := instruction.Parse(to)
		if parsed.GetToken() != instr.GetToken() {
			fmt.Println("Error", parsed.GetToken(), instr.GetToken(), instr.ToASM())
		}
	}
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
	testToAndParse(instrs)
	fmt.Println("Asm:")
	outputASMs(instrs)
	fmt.Println()
	fmt.Println("Bit:")
	outputBits(instrs)

	fmt.Println()
	fmt.Println("Emulate")
	fmt.Println("Initialize:", emulator.Initialize(toBin(instrs)))
	fmt.Println("Execute:", emulator.Execute(0))
	fmt.Println("Registers")
	emulator.ShowRegisters()
}
