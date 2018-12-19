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

type Segment struct {
	Start uint32
	End   uint32
}

type AssembleResult struct {
	Full Segment
	Data Segment
	Text Segment
	Bin  []uint8
}

func assembleWithError(content []string, config AssembleConfig, size int32) (retinstrs []instruction.Instruction, asresult AssembleResult, err error) {

	defer func() {
		if cr := recover(); cr != nil {
			err = cr.(error)
		} else {
			err = nil
		}
	}()

	dataEnd := config.Data
	textEnd := config.Text

	realSize := uint32(0)
	if size > 0 {
		realSize = uint32(size)
	}

	segs := TrimSplitSegment(content)

	buildBits := size > 0

	var result []uint8
	if buildBits {
		result = make([]uint8, size, size)
	} else {
		result = make([]uint8, 0, 0)
	}

	retinstrs = make([]instruction.Instruction, 0)

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
		dataEnd = config.Data + uint32(len(dataSeg))
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
		textEnd = config.Text + (uint32(len(instrs)) << 2)
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
	asresult = AssembleResult{Segment{0, realSize}, Segment{config.Data, dataEnd}, Segment{config.Text, textEnd}, result}
	return retinstrs, asresult, err
}

func Assemble(content []string, config AssembleConfig, size int32) ([]instruction.Instruction, AssembleResult, error) {
	instrs, result, rerr := assembleWithError(content, config, size)
	return instrs, result, rerr
}
