package main

import (
	"fmt"
	"strconv"

	"./emulator"
	ins "./instruction"
)

func createTestInstructions() []ins.Instruction {
	instrs := make([]ins.Instruction, 0, 1024)
	instrs = append(instrs, ins.Lw(2, 0, 0x0))
	instrs = append(instrs, ins.Lw(3, 0, 0x1))
	instrs = append(instrs, ins.Add(1, 2, 3))
	return instrs
}

func toBin(instrs []ins.Instruction) []uint32 {
	result := make([]uint32, 0, len(instrs))
	for _, instr := range instrs {
		result = append(result, instr.ToBits())
	}
	return result
}

func testToAndParse(instrs []ins.Instruction) {
	for _, instr := range instrs {
		to := instr.ToBits()
		parsed := ins.Parse(to)
		if parsed.GetToken() != instr.GetToken() {
			fmt.Println("Error", parsed.GetToken(), instr.GetToken(), instr.ToASM())
		}
	}
}

func outputASMs(instrs []ins.Instruction) {
	for _, instr := range instrs {
		fmt.Println(instr.ToASM())
	}
}

func outputBits(instrs []ins.Instruction) {
	for _, instr := range instrs {
		fmt.Printf("%08s\n", strconv.FormatUint(uint64(instr.ToBits()), 16))
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
