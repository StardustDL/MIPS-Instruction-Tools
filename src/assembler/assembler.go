package assembler

import (
	"fmt"
	"strings"

	"../instruction"
)

type AssembleConfig struct {
	Data uint32
	Text uint32
}

func getInstruction(text string) instruction.Instruction {
	return instruction.CreateJ("", 0, 0)
}

func Assemble(content []string, config AssembleConfig) ([]instruction.Instruction, bool) {
	segs := TrimSplitSegment(content)
	symbolTable := make(map[string]uint32)

	texts, _ := segs["text"]
	flg := true
	result := make([]instruction.Instruction, 0)
	currentAddr := uint32(config.Text)
	for _, str := range texts {
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
	for _, str := range texts {
		if strings.HasSuffix(str, ":") {
			continue
		}
		symbol, tokens := GetTokens(str)
		res, ok := Parse(symbol, tokens, symbolRes, currentAddr+4)
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
