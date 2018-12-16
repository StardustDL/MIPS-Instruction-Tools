package assembler

import (
	// "strings"
	"../instruction"
)

func getInstruction(text string) instruction.Instruction {
	return instruction.CreateJ("", 0, 0)
}

