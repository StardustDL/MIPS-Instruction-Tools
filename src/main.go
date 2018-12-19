package main

import (
	"fmt"
	"strconv"

	ass "./assembler"
	sim "./simulator"
	"./simulator/cpu"
	ins "./instruction"
)

func toASMs(instrs []ins.Instruction) []string {
	result := make([]string, len(instrs))
	for i, instr := range instrs {
		result[i] = instr.ToASM()
	}
	return result
}

func toBitStrings(bin []uint32) []string {
	result := make([]string, len(bin))
	for i, bits := range bin {
		result[i] = fmt.Sprintf("%08s", strconv.FormatUint(uint64(bits), 16))
	}
	return result
}

func toMIF(bin []uint8) []string {
	result := make([]string, 0, len(bin))
	result = append(result, "WIDTH=8;")
	result = append(result, fmt.Sprintf("DEPTH=%d;", len(bin)))
	result = append(result, "ADDRESS_RADIX=HEX;")
	result = append(result, "DATA_RADIX=HEX;")
	result = append(result, "CONTENT BEGIN")
	for i, bits := range bin {
		result = append(result, fmt.Sprintf("    %04X : %04X;", i, bits))
	}
	result = append(result, "END;")
	return result
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
	ENTRY := uint32(0x00001000)
	instrs, bin, ok := ass.Assemble(content, ass.AssembleConfig{Data: 0x00003000, Text: ENTRY}, 0x4000)
	if ok {
		println("done")
	} else {
		println("failed")
		return
	}

	println("Instruction count:", len(instrs))

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

	err = writeAllBytes(OUT_BIN, bin)
	if err != nil {
		println("Generate asm file failed", err)
		return
	}
	println("Bin file:", OUT_BIN)

	err = writeAllLines(OUT_MIF, toMIF(bin))
	if err != nil {
		println("Generate mif file failed", err)
		return
	}
	println("MIF file:", OUT_MIF)

	if len(bin) == 0 {
		println("No bin data. Stop emulating.")
		return
	}

	print("Initializing for simulating...")
	if !sim.Initialize(bin, breakHandler, nil) {
		println("failed")
		return
	}
	println("done")

	println("Executing...")
	flg := sim.Execute(ENTRY, false)
	println("Executed", flg)
	fmt.Println("Registers")
	sim.ShowRegisters()
}

func main() {
	testAssemble()
}
