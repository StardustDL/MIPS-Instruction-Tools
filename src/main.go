package main

import (
    "fmt"
    "strconv"

    ass "./assembler"
    emu "./emulator"
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

const OUT_ASM = "./test/output.asm"
const OUT_BINS = "./test/output.txt"
const OUT_BIN = "./test/output.bin"

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
    instrs, bin, ok := ass.Assemble(content, ass.AssembleConfig{Data: 0x00001800, Text: ENTRY}, 0x2000)
    if ok {
        println("done")
    } else {
        println("failed")
        return
    }

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

    if len(bin) == 0{
        println("No bin data. Stop emulating.")
        return
    }

    print("Initializing for emulating...")
    if !emu.Initialize(bin) {
        println("failed")
        return
    }
    println("done")

    println("Executing...")
    flg := emu.Execute(ENTRY, false)
    println("Executed", flg)
    fmt.Println("Registers")
    emu.ShowRegisters()
}

func main() {
    testAssemble()
}
