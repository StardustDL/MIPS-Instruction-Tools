package assembler

import (
	"strconv"
	"strings"
	"unicode"

	"../instruction"
)

const (
	TC_NONE = uint8(iota)
	TC_REG
	TC_IMM
	TC_SYMBOL
)

type Token struct {
	class  uint8
	value  uint32
	symbol string
}

const regNames = "ATV0V1A0A1A2A3T0T1T2T3T4T5T6T7S0S1S2S3S4S5S6S7T8T9K0K1GPSPFPRA"

func getRegisterToken(val string) Token {
	name := val[1:]
	id := uint32(0)
	if unicode.IsDigit(rune(name[0])) {
		to, _ := strconv.ParseUint(name, 0, 8)
		id = uint32(to)
	} else {
		id = uint32(strings.Index(regNames, name)/2 + 1)
	}
	return Token{TC_REG, id, val}
}

func getImmToken(val string) Token {
	to, _ := strconv.ParseInt(val, 0, 32)
	id := uint32(to)
	return Token{TC_IMM, id, val}
}

func GetTokens(content string) (string, []Token) {
	content = trimLine(content)
	if len(content) == 0 {
		return "", make([]Token, 0)
	}
	indexSpace := strings.Index(content, " ")
	if indexSpace == -1 {
		return content, make([]Token, 0)
	}
	symbol, rem := content[0:indexSpace], content[indexSpace+1:]
	args := strings.Split(rem, ",")
	tokens := make([]Token, 0)
	for _, val := range args {
		val = strings.Trim(val, " \t")
		if strings.HasPrefix(val, "$") { // $r
			tokens = append(tokens, getRegisterToken(val))
		} else if strings.HasSuffix(val, ")") { // offset($r)
			as := strings.Split(val[0:len(val)-1], "(")
			tokens = append(tokens, getRegisterToken(as[1]))
			tokens = append(tokens, getImmToken(as[0]))
		} else { // Imm
			tokens = append(tokens, getImmToken(val))
		}
	}
	return strings.ToLower(symbol), tokens
}

func assertRRR(args []Token) bool {
	return len(args) == 3 && args[0].class == TC_REG && args[1].class == TC_REG && args[2].class == TC_REG
}

func assertRRI(args []Token) bool {
	return len(args) == 3 && args[0].class == TC_REG && args[1].class == TC_REG && args[2].class == TC_IMM
}

func assertRI(args []Token) bool {
	return len(args) == 2 && args[0].class == TC_REG && args[1].class == TC_IMM
}

func assertI(args []Token) bool {
	return len(args) == 1 && args[0].class == TC_IMM
}

func assertR(args []Token) bool {
	return len(args) == 1 && args[0].class == TC_REG
}

type SymbolResolver func(args []Token) []Token

func Parse(symbol string, args []Token, resolver SymbolResolver) (instruction.Instruction, bool) {
	args = resolver(args)
	switch symbol {
	case "add":
		if assertRRR(args) {
			break
		}
		return instruction.Add(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "addu":
		if assertRRR(args) {
			break
		}
		return instruction.Addu(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "addi":
		if assertRRI(args) {
			break
		}
		return instruction.Addi(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "addiu":
		if assertRRI(args) {
			break
		}
		return instruction.Addiu(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sub":
		if assertRRR(args) {
			break
		}
		return instruction.Sub(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "subu":
		if assertRRR(args) {
			break
		}
		return instruction.Subu(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "and":
		if assertRRR(args) {
			break
		}
		return instruction.And(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "andi":
		if assertRRI(args) {
			break
		}
		return instruction.Andi(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "or":
		if assertRRR(args) {
			break
		}
		return instruction.Or(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "ori":
		if assertRRI(args) {
			break
		}
		return instruction.Ori(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "nor":
		if assertRRR(args) {
			break
		}
		return instruction.Nor(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "xor":
		if assertRRR(args) {
			break
		}
		return instruction.Xor(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "xori":
		if assertRRI(args) {
			break
		}
		return instruction.Xori(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "slt":
		if assertRRR(args) {
			break
		}
		return instruction.Slt(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "slti":
		if assertRRI(args) {
			break
		}
		return instruction.Slti(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sltu":
		if assertRRR(args) {
			break
		}
		return instruction.Sltu(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "sltiu":
		if assertRRI(args) {
			break
		}
		return instruction.Sltiu(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sll":
		if assertRRI(args) {
			break
		}
		return instruction.Sll(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "sra":
		if assertRRI(args) {
			break
		}
		return instruction.Sra(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "srl":
		if assertRRI(args) {
			break
		}
		return instruction.Srl(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "sllv":
		if assertRRR(args) {
			break
		}
		return instruction.Sllv(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "srav":
		if assertRRR(args) {
			break
		}
		return instruction.Srav(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "srlv":
		if assertRRR(args) {
			break
		}
		return instruction.Srlv(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "beq":
		if assertRRI(args) {
			break
		}
		return instruction.Beq(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "bne":
		if assertRRI(args) {
			break
		}
		return instruction.Bne(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "bgez":
		if assertRI(args) {
			break
		}
		return instruction.Bgez(uint8(args[0].value), uint16(args[1].value)), true
	case "bgezal":
		if assertRI(args) {
			break
		}
		return instruction.Bgezal(uint8(args[0].value), uint16(args[1].value)), true
	case "bgtz":
		if assertRI(args) {
			break
		}
		return instruction.Bgtz(uint8(args[0].value), uint16(args[1].value)), true
	case "blez":
		if assertRI(args) {
			break
		}
		return instruction.Blez(uint8(args[0].value), uint16(args[1].value)), true
	case "bltz":
		if assertRI(args) {
			break
		}
		return instruction.Bltz(uint8(args[0].value), uint16(args[1].value)), true
	case "bltzal":
		if assertRI(args) {
			break
		}
		return instruction.Bltzal(uint8(args[0].value), uint16(args[1].value)), true
	case "j":
		if assertI(args) {
			break
		}
		return instruction.J(uint32(args[0].value)), true
	case "jal":
		if assertI(args) {
			break
		}
		return instruction.Jal(uint32(args[0].value)), true
	case "jr":
		if assertR(args) {
			break
		}
		return instruction.Jr(uint8(args[0].value)), true
	case "lui":
		if assertRI(args) {
			break
		}
		return instruction.Lui(uint8(args[0].value), uint16(args[1].value)), true
	case "lb":
		if assertRRI(args) {
			break
		}
		return instruction.Lb(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sb":
		if assertRRI(args) {
			break
		}
		return instruction.Sb(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "lw":
		if assertRRI(args) {
			break
		}
		return instruction.Lw(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sw":
		if assertRRI(args) {
			break
		}
		return instruction.Sw(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	}
	return instruction.Nop(), false
}
