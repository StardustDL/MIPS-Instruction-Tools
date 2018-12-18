package main

import (
    "fmt"
    ins "./instruction"
)

func createTestInstructions() []ins.Instruction {
    instrs := make([]ins.Instruction, 0, 1024)
    instrs = append(instrs, ins.Lw(2, 0, 0x0))
    instrs = append(instrs, ins.Lw(3, 0, 0x1))
    instrs = append(instrs, ins.Add(1, 2, 3))
    return instrs
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