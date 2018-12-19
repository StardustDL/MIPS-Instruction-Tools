package assembler

import (
    "fmt"
    // "strings"

    "../instruction"
)

type AssembleConfig struct {
    Data uint32
    Text uint32
}

func Assemble(content []string, config AssembleConfig, size int32) ([]instruction.Instruction, []uint8, bool) {
    segs := TrimSplitSegment(content)

    buildBits := size > 0

    var result []uint8
    if buildBits {
        result = make([]uint8, size, size)
    } else {
        result = make([]uint8, 0, 0)
    }

    retinstrs := make([]instruction.Instruction, 0)

    symbolTable := make(map[string]uint32)

    flg := true

    data, hasData := segs["data"]
    if hasData {
        dataSeg, dataSymbols, ok := buildData(data, config, buildBits)
        if !ok {
            fmt.Printf("Data segment failed.")
            flg = false
        }
        for k, v := range dataSymbols {
            _, exists := symbolTable[k]
            if exists {
                flg = false
                fmt.Printf("Symbol %s has been defined.\n", k)
                break
            }
            symbolTable[k] = v
        }
        if buildBits {
            for i, val := range dataSeg {
                result[config.Data+uint32(i)] = val
            }
        }
    }

    texts, hasText := segs["text"]
    if flg && hasText {
        instrs, ok := buildText(texts, config, symbolTable)
        if !ok {
            fmt.Printf("Text segment failed.")
            flg = false
        }
        retinstrs = instrs
        if buildBits {
            for i, bits := range instruction.ToBin(instrs) {
                for j := uint32(0); j < 4; j++ {
                    result[config.Text+(uint32(i)<<2)+j] = uint8(bits >> (j << 3) & 0xff)
                }
            }
        }
    }

    if len(segs[DEFAULT_SEGMENT]) > 0 {
        println("Warning: some instruction not in any special segment")
    }

    return retinstrs, result, flg
}
