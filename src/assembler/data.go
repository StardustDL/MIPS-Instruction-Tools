package assembler

import (
	"fmt"
	"regexp"
	"strings"
	"errors"
)

var dataTokenRegex *regexp.Regexp
var groupNames []string

func getDataTokens(str string) (string, string, []uint8) {
	match := dataTokenRegex.FindStringSubmatch(str)
	result := make(map[string]string)
	if len(match)<len(groupNames) {
		panic(errors.New("Data token parsing failed"))
	}
	for i, name := range groupNames {
		if i != 0 && name != "" { // 第一个分组为空（也就是整个匹配）
			result[name] = match[i]
		}
	}
	data := make([]uint8, 0)
	switch result["type"] {
	case "asciiz":
		raw := strings.Trim(result["content"], "\"")
		raw = strings.Replace(raw, "\\n", "\n", -1) // TODO
		raw = strings.Replace(raw, "\\r", "\r", -1)
		for _, chr := range raw {
			data = append(data, uint8(chr))
		}
		data = append(data, uint8(0))
	}
	return result["symbol"], result["type"], data
}

func buildData(content []string, config AssembleConfig, buildBits bool) ([]uint8, map[string]uint32) {
	dataTokenRegex = regexp.MustCompile(`^(?P<symbol>[\w]+)[\s]*:[\s]*\.(?P<type>[\w]+)[\s]*(?P<content>\S[\s\S]*)$`)
	groupNames = dataTokenRegex.SubexpNames()
	result := make([]uint8, 0)
	symbolTable := make(map[string]uint32)
	dataOffset := config.Data
	for _, str := range content {
		symbol, _, data := getDataTokens(str)
		_, exists := symbolTable[symbol]
		if exists {
			panic(errors.New(fmt.Sprintf("Symbol %s has been defined.", symbol)))
		}
		symbolTable[symbol] = dataOffset
		if buildBits {
			for _, val := range data {
				result = append(result, val)
				dataOffset++
			}
		} else {
			dataOffset += uint32(len(data))
		}
	}
	return result, symbolTable
}
