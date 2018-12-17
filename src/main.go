package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	ass "./assembler"
	emu "./emulator"
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

func outputASMs(instrs []ins.Instruction) {
	for _, instr := range instrs {
		fmt.Println(instr.ToASM())
	}
}

func outputBits(instrs []ins.Instruction) {
	for _, bits := range ins.ToBin(instrs) {
		fmt.Printf("%08s\n", strconv.FormatUint(uint64(bits), 16))
	}
}

const (
	PORT_INPUT  = 0x0000
	PORT_OUTPUT = 0x0080
	PORT_LED    = 0x0100
	PORT_DIGIT  = 0x0101
	SEG_TEXT    = 0x1000
	ENTRY       = 0x1000
)

var memory [emu.MEMORY_SIZE]uint8

func readAllLines(path string) []string {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return make([]string, 0)
	}
	defer file.Close()

	_, err = file.Stat()
	if err != nil {
		panic(err)
	}

	buf := bufio.NewReader(file)

	result := make([]string, 0)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		result = append(result, line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				break
			}
		}
	}
	return result
}

func testAssemble() {
	content := readAllLines("./test.S")
	ENTRY := uint32(0x00001000)
	// instrs, _, ok := ass.Assemble(content, ass.AssembleConfig{Data: 0x10010000, Text: 0x00400000, EndInstruction: emu.END_INSTR}, -1)
	instrs, bin, ok := ass.Assemble(content, ass.AssembleConfig{Data: 0x00001800, Text: ENTRY}, 0x2000)
	if !ok {
		fmt.Println("Assembling failed.")
		return
	}
	fmt.Println("Asm:")
	outputASMs(instrs)
	fmt.Println()
	// fmt.Println("Bit:")
	// outputBits(instrs)
	fmt.Println()
	fmt.Println("Emulate")
	fmt.Println("Initializing")
	if !emu.Initialize(bin) {
		fmt.Println("Initializing failed")
	} else {
		fmt.Println("Executing")
		fmt.Println("Executed:", emu.Execute(ENTRY))
		fmt.Println("Registers")
		emu.ShowRegisters()
	}
}

func main() {
	testAssemble()
}
