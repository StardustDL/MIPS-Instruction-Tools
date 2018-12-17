package assembler

import (
	"fmt"
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

type InstructionSyntax struct {
	symbol string
	args   []Token
}

const regNames = "atv0v1a0a1a2a3t0t1t2t3t4t5t6t7s0s1s2s3s4s5s6s7t8t9k0k1gpspfpra"

func getRegisterToken(val string) Token {
	name := strings.ToLower(val[1:])
	id := uint32(0)
	if unicode.IsDigit(rune(name[0])) {
		to, _ := strconv.ParseUint(name, 0, 8)
		id = uint32(to)
	} else if name == "zero" {
		id = uint32(0)
	} else {
		ind := strings.Index(regNames, name)
		if ind != -1 {
			id = uint32(ind/2 + 1)
		} else {
			fmt.Printf("No this register: %s\n", name)
		}
	}
	return Token{TC_REG, id, val}
}

func getImmOrSymToken(val string) Token {
	to, err := strconv.ParseInt(val, 0, 32)
	if err == nil {
		id := uint32(to)
		return Token{TC_IMM, id, val}
	} else if len(val) == 0 {
		return Token{TC_IMM, 0, val}
	} else {
		return Token{TC_SYMBOL, 0, val}
	}
}

func getTextTokens(content string) InstructionSyntax {
	content = trimLine(content)
	if len(content) == 0 {
		return InstructionSyntax{"", make([]Token, 0)}
	}
	indexSpace := strings.Index(content, " ")
	if indexSpace == -1 {
		return InstructionSyntax{content, make([]Token, 0)}
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
			tokens = append(tokens, getImmOrSymToken(as[0]))
		} else { // Imm
			tokens = append(tokens, getImmOrSymToken(val))
		}
	}
	return InstructionSyntax{strings.ToLower(symbol), tokens}
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

func TextPreprocess(content []string, resolver SymbolResolver, config AssembleConfig) ([]InstructionSyntax, map[string]uint32, bool) {
	currentAddr := uint32(config.Text)
	syntaxs := make([]InstructionSyntax, 0, len(content))
	symbolTable := make(map[string]uint32)
	flg := true
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
		} else {
			syntax := getTextTokens(str)
			syntax.args = resolver(syntax.args)
			tosyn, ok := textPreprocessOne(syntax)
			if !ok {
				flg = false
				fmt.Printf("Preprocessing failed: %s\n", str)
				break
			}
			for _, v := range tosyn {
				syntaxs = append(syntaxs, v)
				currentAddr += 4
			}
		}
	}
	return syntaxs, symbolTable, flg
}

func textPreprocessOne(syntax InstructionSyntax) ([]InstructionSyntax, bool) {
	switch syntax.symbol {
	case "li": // extra instr
		if !assertRI(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"lui", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_AT)},
				Token{class: TC_IMM, value: syntax.args[1].value >> 16}}},
			InstructionSyntax{"ori", []Token{
				Token{class: TC_REG, value: uint32(syntax.args[0].value)},
				Token{class: TC_REG, value: uint32(instruction.GPR_AT)},
				Token{class: TC_IMM, value: syntax.args[1].value & 0xffff}}}}, true
	case "":
		break
	default:
		return []InstructionSyntax{syntax}, true

	}
	return []InstructionSyntax{syntax}, false
}

func TextParse(syntax InstructionSyntax, resolver SymbolResolver, nextPC uint32) (instruction.Instruction, bool) {
	args := resolver(syntax.args)
	switch syntax.symbol {
	case "add":
		if !assertRRR(args) {
			break
		}
		return instruction.Add(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "addu":
		if !assertRRR(args) {
			break
		}
		return instruction.Addu(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "addi":
		if !assertRRI(args) {
			break
		}
		return instruction.Addi(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "addiu":
		if !assertRRI(args) {
			break
		}
		return instruction.Addiu(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sub":
		if !assertRRR(args) {
			break
		}
		return instruction.Sub(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "subu":
		if !assertRRR(args) {
			break
		}
		return instruction.Subu(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "and":
		if !assertRRR(args) {
			break
		}
		return instruction.And(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "andi":
		if !assertRRI(args) {
			break
		}
		return instruction.Andi(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "or":
		if !assertRRR(args) {
			break
		}
		return instruction.Or(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "ori":
		if !assertRRI(args) {
			break
		}
		return instruction.Ori(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "nor":
		if !assertRRR(args) {
			break
		}
		return instruction.Nor(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "xor":
		if !assertRRR(args) {
			break
		}
		return instruction.Xor(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "xori":
		if !assertRRI(args) {
			break
		}
		return instruction.Xori(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "slt":
		if !assertRRR(args) {
			break
		}
		return instruction.Slt(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "slti":
		if !assertRRI(args) {
			break
		}
		return instruction.Slti(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sltu":
		if !assertRRR(args) {
			break
		}
		return instruction.Sltu(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "sltiu":
		if !assertRRI(args) {
			break
		}
		return instruction.Sltiu(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
	case "sll":
		if !assertRRI(args) {
			break
		}
		return instruction.Sll(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "sra":
		if !assertRRI(args) {
			break
		}
		return instruction.Sra(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "srl":
		if !assertRRI(args) {
			break
		}
		return instruction.Srl(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "sllv":
		if !assertRRR(args) {
			break
		}
		return instruction.Sllv(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "srav":
		if !assertRRR(args) {
			break
		}
		return instruction.Srav(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "srlv":
		if !assertRRR(args) {
			break
		}
		return instruction.Srlv(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "beq":
		if !assertRRI(args) {
			break
		}
		return instruction.Beq(uint8(args[0].value), uint8(args[1].value), uint16((args[2].value-nextPC)>>2)), true
	case "bne":
		if !assertRRI(args) {
			break
		}
		return instruction.Bne(uint8(args[0].value), uint8(args[1].value), uint16((args[2].value-nextPC)>>2)), true
	case "bgez":
		if !assertRI(args) {
			break
		}
		return instruction.Bgez(uint8(args[0].value), uint16((args[1].value-nextPC)>>2)), true
	case "bgezal":
		if !assertRI(args) {
			break
		}
		return instruction.Bgezal(uint8(args[0].value), uint16((args[1].value-nextPC)>>2)), true
	case "bgtz":
		if !assertRI(args) {
			break
		}
		return instruction.Bgtz(uint8(args[0].value), uint16((args[1].value-nextPC)>>2)), true
	case "blez":
		if !assertRI(args) {
			break
		}
		return instruction.Blez(uint8(args[0].value), uint16((args[1].value-nextPC)>>2)), true
	case "bltz":
		if !assertRI(args) {
			break
		}
		return instruction.Bltz(uint8(args[0].value), uint16((args[1].value-nextPC)>>2)), true
	case "bltzal":
		if !assertRI(args) {
			break
		}
		return instruction.Bltzal(uint8(args[0].value), uint16((args[1].value-nextPC)>>2)), true
	case "j":
		if !assertI(args) {
			break
		}
		return instruction.J(uint32(args[0].value)), true
	case "jal":
		if !assertI(args) {
			break
		}
		return instruction.Jal(uint32(args[0].value)), true
	case "jr":
		if !assertR(args) {
			break
		}
		return instruction.Jr(uint8(args[0].value)), true
	case "lui":
		if !assertRI(args) {
			break
		}
		return instruction.Lui(uint8(args[0].value), uint16(args[1].value)), true
	case "lb":
		switch len(args) {
		case 3:
			if !assertRRI(args) {
				break
			}
			return instruction.Lb(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		case 2:
			if !assertRI(args) {
				break
			}
			return instruction.Lb(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "sb":
		switch len(args) {
		case 3:
			if !assertRRI(args) {
				break
			}
			return instruction.Sb(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		case 2:
			if !assertRI(args) {
				break
			}
			return instruction.Sb(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "lw":
		switch len(args) {
		case 3:
			if !assertRRI(args) {
				break
			}
			return instruction.Lw(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		case 2:
			if !assertRI(args) {
				break
			}
			return instruction.Lw(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "sw":
		switch len(args) {
		case 3:
			if !assertRRI(args) {
				break
			}
			return instruction.Sw(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		case 2:
			if !assertRI(args) {
				break
			}
			return instruction.Sw(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "nop":
		if len(args) > 0 {
			break
		}
		return instruction.Nop(), true
	case "syscall":
		if len(args) > 0 {
			break
		}
		return instruction.Syscall(), true
	}
	return instruction.Nop(), false
}
