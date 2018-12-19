package assembler

import (
	"errors"
	"fmt"
	// "strings"

	"../instruction"
)

type AssembleConfig struct {
	Data uint32
	Text uint32
}

func assembleWithError(content []string, config AssembleConfig, size int32) ([]instruction.Instruction, []uint8, *error) {
	var err error

	defer func() {
		if cr := recover(); cr != nil {
			err = cr.(error)
		} else {
			err = nil
		}
	}()

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

	data, hasData := segs["data"]
	if hasData {
		dataSeg, dataSymbols := buildData(data, config, buildBits)
		for k, v := range dataSymbols {
			_, exists := symbolTable[k]
			if exists {
				panic(errors.New(fmt.Sprintf("Symbol %s has been defined.", k)))
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
	if hasText {
		instrs := buildText(texts, config, symbolTable)
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

	return retinstrs, result, &err
}

func Assemble(content []string, config AssembleConfig, size int32) ([]instruction.Instruction, []uint8, error) {
	instrs, bin, err := assembleWithError(content, config, size)
	return instrs, bin, *err
}
