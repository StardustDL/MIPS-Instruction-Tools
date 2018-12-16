package main

import (
	"fmt"
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
	for j := uint32(0); j < 4; j++ {
		memory[SEG_TEXT+(uint32(len(text))<<2)+j] = uint8(emu.END_INSTR >> (j << 3) & 0xff)
	}

	emu.Initialize(memory[:])
}

func testAssemble() {
	content := `
	.text

main:
jal read
nop
lui $a1,0x00000000
lui $a2,0x00001001
addiu $a1,$a1,0x1800
addiu $a2,$a2,0
jal streq
nop
bne $v0,$zero,hello_world
nop
hello_world:
lui $a1,0x00000000
addiu $a1,$a1,0x1800
jal write
nop
j main
nop

streq: # compare string s1 and s2
streq_loop:
lb $t1,($a1)
lb $t2,($a2)
or $t3,$t1,$t2
beq $t3,$zero,streq_eq
nop
bne $t1,$t2,streq_ne
nop
addiu $a1,$a1,1
addiu $a2,$a2,1
j streq_loop
nop
streq_ne: 
addiu $v0,$zero,0
jr $ra
nop
streq_eq:
addiu $v0,$zero,1
jr $ra
nop

strcpy:
strcpy_loop:
lb $t1,($a1)
sb $t1,($a2)
beq $t1,$zero,strcpy_end
nop
addiu $a1,$a1,1
addiu $a2,$a2,1
j strcpy_loop
nop
strcpy_end:
jr $ra
nop

read:
read_before:
lb $t1,0x0000
beq $t1,$zero,read_before
nop
j read_main
nop
read_main:
lui $a1,0x0000
lui $a2,0x1800
jal strcpy
nop
sb $zero,0x0000
jr $ra
nop

write:
lui $a2,0x0080
jal strcpy
nop
jr $ra
nop
	`
	result, ok := ass.Assemble(strings.Split(content, "\n"), ass.AssembleConfig{Data: 0x0, Text: 0x00400000})
	if !ok {
		fmt.Println("Assembling failed.")
		return
	}
	fmt.Println("Asm:")
	outputASMs(result)
	fmt.Println()
	fmt.Println("Bit:")
	outputBits(result)
}

func main() {
	testAssemble()
	return
	instrs = createTestInstructions()
	testToAndParse(instrs)
	fmt.Println("Asm:")
	outputASMs(instrs)
	fmt.Println()
	fmt.Println("Bit:")
	outputBits(instrs)

	fmt.Println()
	fmt.Println("Emulate")
	fmt.Println("Initializing")
	initMemory(ins.ToBin(instrs))
	fmt.Println("Executing")
	fmt.Println("Executed:", emu.Execute(ENTRY))
	fmt.Println("Registers")
	emu.ShowRegisters()
}
