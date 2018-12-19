package main

import (
    "fmt"

	ass "./assembler"
	ins "./instruction"
	sim "./simulator"
	"./simulator/cpu"
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

const OUT_ASM = "./test/output.asm"
const OUT_BINS = "./test/output.txt"
const OUT_BIN = "./test/output.bin"
const OUT_MIF = "./test/output.mif"

func breakHandler(code uint32) {
	// just for debug

	println(fmt.Sprintf("Recieve break code: %d", code))

	if code == 100 {
		n := cpu.GetGPR(ins.GPR_A0)
		st := cpu.GetGPR(ins.GPR_A1)
		println("LED", n-2, "turn to", st)
	} else if code == 200 {
		println("Clear console")
	} else if code == 1 {
		a0 := cpu.GetGPR(ins.GPR_A0)
		a1 := cpu.GetGPR(ins.GPR_A1)
		println("a0:", a0, "a1:", a1)
	}
}

func testAssemble() {
	content, err := readAllLines("./test/input.asm")
	if err != nil {
		println("File reading error", err)
		return
	}

	print("Assembling...")

	// ENTRY := uint32(0x00400000)
	// instrs, bin, ok := ass.Assemble(content, ass.AssembleConfig{Data: 0x10010000, Text: ENTRY}, -1)
	instrs, builded, err := ass.Assemble(content, ass.AssembleConfig{Data: 0x00003000, Text: 0x00001000}, 0x4000)
	if err == nil {
		println("done")
	} else {
		println("failed")
		println(err.Error())
		return
	}

	println("Instruction count:", len(instrs))
	fmt.Printf("Full segment: 0x%08x ~ 0x%08x\n", builded.Full.Start, builded.Full.End)
	fmt.Printf("Data segment: 0x%08x ~ 0x%08x\n", builded.Data.Start, builded.Data.End)
	fmt.Printf("Text segment: 0x%08x ~ 0x%08x\n", builded.Text.Start, builded.Text.End)

	err = writeAllLines(OUT_BINS, toBitStrings(ins.ToBin(instrs)))
	if err != nil {
		println("Generate bit string file failed", err)
		return
	}
	println("Bit string file:", OUT_BINS)

	err = writeAllLines(OUT_ASM, toASMs(instrs))
	if err != nil {
		println("Generate asm file failed", err)
		return
	}
	println("ASM file:", OUT_ASM)

	err = writeAllBytes(OUT_BIN, builded.Bin)
	if err != nil {
		println("Generate asm file failed", err)
		return
	}
	println("Bin file:", OUT_BIN)

	err = writeAllLines(OUT_MIF, toMIF(builded.Bin))
	if err != nil {
		println("Generate mif file failed", err)
		return
	}
	println("MIF file:", OUT_MIF)

	if builded.Full.End == 0 {
		println("No bin data. Stop emulating.")
		return
	}

	print("Initializing for simulating...")
	if !sim.Initialize(builded.Bin, breakHandler, nil) {
		println("failed")
		return
	}
	println("done")

	println("Executing...")
	flg := sim.Execute(builded.Text.Start, false)
	println("Executed", flg)
	fmt.Println("Registers")
	sim.ShowRegisters()
}
