package assembler

import (
	"errors"
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
		if len(name)==2 && ind != -1 && ind % 2 == 0{
			id = uint32(ind/2 + 1)
		} else {
			panic(errors.New(fmt.Sprintf("No this register: %s", name)))
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

func assertRRRn(args []Token) bool {
	return len(args) == 3 && args[0].class == TC_REG && args[1].class == TC_REG && args[2].class == TC_REG
}

func assertRI(args []Token) bool {
	return len(args) == 2 && args[0].class == TC_REG && args[1].class == TC_IMM
}

func assertRRn(args []Token) bool {
	return len(args) == 2 && args[0].class == TC_REG && args[1].class != TC_REG
}

func assertRR(args []Token) bool {
	return len(args) == 2 && args[0].class == TC_REG && args[1].class == TC_REG
}

func assertI(args []Token) bool {
	return len(args) == 1 && args[0].class == TC_IMM
}

func assertR(args []Token) bool {
	return len(args) == 1 && args[0].class == TC_REG
}

func assertRn(args []Token) bool {
	return len(args) == 1 && args[0].class != TC_REG
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
	tokenZERO := Token{class: TC_REG, value: uint32(instruction.GPR_ZERO)}
	tokenAT := Token{class: TC_REG, value: uint32(instruction.GPR_AT)}

	switch syntax.symbol {
	case "li":
		if !assertRRn(syntax.args) {
			break
		}
		imm := syntax.args[1].value
		if imm <= 0xffff { // dont have high 16 bit
			if (imm & 0x8000) > 0 { // sign bit is 1
				return []InstructionSyntax{
					InstructionSyntax{"ori", []Token{
						Token{class: TC_REG, value: uint32(syntax.args[0].value)},
						tokenAT,
						Token{class: TC_IMM, value: imm}}},
				}, true
			} else {
				return []InstructionSyntax{
					InstructionSyntax{"addiu", []Token{
						Token{class: TC_REG, value: uint32(syntax.args[0].value)},
						tokenZERO,
						Token{class: TC_IMM, value: imm}}},
				}, true
			}
		} else {
			return []InstructionSyntax{
				InstructionSyntax{"lui", []Token{
					tokenAT,
					Token{class: TC_IMM, value: imm >> 16}}},
				InstructionSyntax{"ori", []Token{
					Token{class: TC_REG, value: uint32(syntax.args[0].value)},
					tokenAT,
					Token{class: TC_IMM, value: imm & 0xffff}}},
			}, true
		}
	case "la":
		if assertRRRn(syntax.args) {
			return []InstructionSyntax{
				InstructionSyntax{"ori", []Token{
					tokenAT,
					tokenZERO,
					syntax.args[2]}},
				InstructionSyntax{"add", []Token{
					syntax.args[0],
					syntax.args[1],
					tokenAT}},
			}, true
		} else if assertRRn(syntax.args) {
			return []InstructionSyntax{
				InstructionSyntax{"addi", []Token{
					syntax.args[0],
					tokenZERO,
					syntax.args[1]}},
			}, true
		}
	case "neg":
		if !assertRR(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"sub", []Token{
				syntax.args[0],
				tokenZERO,
				syntax.args[1]}},
		}, true
	case "negu":
		if !assertRR(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"subu", []Token{
				syntax.args[0],
				tokenZERO,
				syntax.args[1]}},
		}, true
	case "move":
		if !assertRR(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"addu", []Token{
				syntax.args[0],
				tokenZERO,
				syntax.args[1]}},
		}, true
	case "bal":
		if !assertRn(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"bgezal", []Token{
				tokenZERO,
				syntax.args[0]}},
		}, true
	case "b":
		if !assertRn(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"beq", []Token{
				tokenZERO,
				tokenZERO,
				syntax.args[0]}},
		}, true
	case "beqz":
		if !assertRRn(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"beq", []Token{
				syntax.args[0],
				tokenZERO,
				syntax.args[1]}},
		}, true
	case "bnez":
		if !assertRRn(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"bne", []Token{
				syntax.args[0],
				tokenZERO,
				syntax.args[1]}},
		}, true
	case "not":
		if !assertRR(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"nor", []Token{
				syntax.args[0],
				syntax.args[1],
				tokenZERO}}}, true
	case "nop":
		if !(len(syntax.args) == 0) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"sll", []Token{
				tokenZERO,
				tokenZERO,
				Token{class: TC_IMM, value: 0}}},
		}, true
	case "push": // subu $sp,$sp,4; sw $r,($sp)
		if !assertR(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"addiu", []Token{
				tokenAT,
				tokenZERO,
				Token{class: TC_IMM, value: 4}}},
			InstructionSyntax{"subu", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				tokenAT}},
			InstructionSyntax{"sw", []Token{
				syntax.args[0],
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_IMM, value: 0}}},
		}, true
	case "pop": // lw $r,($sp); addiu $sp,$sp,4
		if !assertR(syntax.args) {
			break
		}
		return []InstructionSyntax{
			InstructionSyntax{"lw", []Token{
				syntax.args[0],
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_IMM, value: 0}}},
			InstructionSyntax{"addiu", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_IMM, value: 4}}},
		}, true
	case "call":
		if !assertRn(syntax.args) { // only support imm or symbol
			break
		}
		return []InstructionSyntax{
			// push $ra
			InstructionSyntax{"addiu", []Token{
				tokenAT,
				tokenZERO,
				Token{class: TC_IMM, value: 4}}},
			InstructionSyntax{"subu", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				tokenAT}},
			InstructionSyntax{"sw", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_RA)},
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_IMM, value: 0}}},
			// jal addr
			InstructionSyntax{"jal", []Token{
				syntax.args[0]}},
		}, true
	case "ret":
		if !(len(syntax.args) == 0) {
			break
		}
		return []InstructionSyntax{
			// move $ra, $at
			InstructionSyntax{"addu", []Token{
				tokenAT,
				tokenZERO,
				Token{class: TC_REG, value: uint32(instruction.GPR_RA)}}},
			// pop $ra
			InstructionSyntax{"lw", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_RA)},
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_IMM, value: 0}}},
			InstructionSyntax{"addiu", []Token{
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_REG, value: uint32(instruction.GPR_SP)},
				Token{class: TC_IMM, value: 4}}},
			// jr $at
			InstructionSyntax{"jr", []Token{
				tokenAT,
			}},
		}, true
	case "":
		break
	default:
		return []InstructionSyntax{syntax}, true

	}
	return []InstructionSyntax{syntax}, false
}

func textParseOne(syntax InstructionSyntax, resolver SymbolResolver, nextPC uint32) (instruction.Instruction, bool) {
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
	case "mul":
		if !assertRRR(args) {
			break
		}
		return instruction.Mul(uint8(args[0].value), uint8(args[1].value), uint8(args[2].value)), true
	case "div":
		if !assertRR(args) {
			break
		}
		return instruction.Div(uint8(args[0].value), uint8(args[1].value)), true
	case "divu":
		if !assertRR(args) {
			break
		}
		return instruction.Divu(uint8(args[0].value), uint8(args[1].value)), true
	case "mult":
		if !assertRR(args) {
			break
		}
		return instruction.Mult(uint8(args[0].value), uint8(args[1].value)), true
	case "multu":
		if !assertRR(args) {
			break
		}
		return instruction.Multu(uint8(args[0].value), uint8(args[1].value)), true
	case "mfhi":
		if !assertR(args) {
			break
		}
		return instruction.Mfhi(uint8(args[0].value)), true
	case "mthi":
		if !assertR(args) {
			break
		}
		return instruction.Mthi(uint8(args[0].value)), true
	case "mflo":
		if !assertR(args) {
			break
		}
		return instruction.Mflo(uint8(args[0].value)), true
	case "mtlo":
		if !assertR(args) {
			break
		}
		return instruction.Mtlo(uint8(args[0].value)), true
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
	case "jalr":
		if assertRR(args) {
			return instruction.Jalr(uint8(args[0].value), uint8(args[1].value)), true
		} else if assertR(args) {
			return instruction.Jalr(instruction.GPR_RA, uint8(args[0].value)), true
		}
	case "lui":
		if !assertRI(args) {
			break
		}
		return instruction.Lui(uint8(args[0].value), uint16(args[1].value)), true
	case "lb":
		if assertRRI(args) {
			return instruction.Lb(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Lb(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "lbu":
		if assertRRI(args) {
			return instruction.Lbu(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Lbu(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "lh":
		if assertRRI(args) {
			return instruction.Lh(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Lh(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "lhu":
		if assertRRI(args) {
			return instruction.Lhu(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Lhu(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "sb":
		if assertRRI(args) {
			return instruction.Sb(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Sb(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "lw":
		if assertRRI(args) {
			return instruction.Lw(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Lw(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "sw":
		if assertRRI(args) {
			return instruction.Sw(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Sw(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "sh":
		if assertRRI(args) {
			return instruction.Sh(uint8(args[0].value), uint8(args[1].value), uint16(args[2].value)), true
		} else if assertRI(args) {
			return instruction.Sh(uint8(args[0].value), uint8(0), uint16(args[1].value)), true
		}
	case "syscall":
		if len(args) > 0 {
			break
		}
		return instruction.Syscall(), true
	case "break":
		if len(args) == 0 {
			return instruction.Break(0), true
		} else if assertI(args) {
			return instruction.Break(args[0].value), true
		}
	}
	return nil, false
}

func buildText(content []string, config AssembleConfig, symbolTable map[string]uint32) []instruction.Instruction {
	result := make([]instruction.Instruction, 0)

	symbolResWithoutError := func(args []Token) []Token {
		for i, item := range args {
			if item.class == TC_SYMBOL {
				val, ok := symbolTable[item.symbol]
				if ok {
					args[i] = Token{TC_IMM, val, item.symbol}
				} else {

				}
			}
		}
		return args
	}

	syntaxs, textSymbols, ok := TextPreprocess(content, symbolResWithoutError, config)

	if ok {
		for k, v := range textSymbols {
			_, exists := symbolTable[k]
			if exists {
				panic(errors.New(fmt.Sprintf("Symbol %s has been defined.", k)))
			}
			symbolTable[k] = v
		}
	} else {
		panic(errors.New(fmt.Sprintf("Text prepocessing failed.")))
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

	currentAddr := uint32(config.Text)
	for _, syn := range syntaxs {
		res, ok := textParseOne(syn, symbolRes, currentAddr+4)
		if !ok {
			panic(errors.New(fmt.Sprintf("Parse failed: %s\n", syn.symbol)))
		}
		result = append(result, res)
		currentAddr += 4
	}
	return result
}
