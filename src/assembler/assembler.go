package assembler

import (
	"fmt"
	"strings"

	"../instruction"
)

type AssembleConfig struct {
	Data           uint32
	Text           uint32
	EndInstruction uint32
}

func getInstruction(text string) instruction.Instruction {
	return instruction.CreateJ("", 0, 0)
}

func buildText(content []string, config AssembleConfig, symbolTable map[string]uint32) ([]instruction.Instruction, bool) {
	result := make([]instruction.Instruction, 0)
	flg := true
	currentAddr := uint32(config.Text)
	for _, str := range content {
		if strings.HasSuffix(str, ":") {
			name := str[0 : len(str)-1]
			_, exists := symbolTable[name]
			if exists {
				flg = false
				fmt.Printf("Symbol %s has been defined.\n", name)
				break
			}
			symbolTable[name] = currentAddr
			continue
		}
		currentAddr += 4
	}

	symbolRes := func(args []Token) []Token {
		for i, item := range args {
			if item.class == TC_SYMBOL {
				val, ok := symbolTable[item.symbol]
				if ok {
					args[i] = Token{TC_IMM, val, item.symbol}
				} else {
					fmt.Printf("No this symbol: %s\n", item.symbol)
				}
			}
		}
		return args
	}

	currentAddr = uint32(config.Text)
	for _, str := range content {
		if strings.HasSuffix(str, ":") {
			continue
		}
		symbol, tokens := getTextTokens(str)
		res, ok := TextParse(symbol, tokens, symbolRes, currentAddr+4)
		if !ok {
			flg = false
			fmt.Printf("Parse failed: %s\n", str)
			break
		}
		result = append(result, res)
		currentAddr += 4
	}
	return result, flg
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
			for j := uint32(0); j < 4; j++ {
				result[config.Text+(uint32(len(instrs))<<2)+j] = uint8(config.EndInstruction >> (j << 3) & 0xff)
			}
		}
	}

	return retinstrs, result, flg
}
