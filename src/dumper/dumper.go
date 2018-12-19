package dumper

import (
	ins "../instruction"
)

func DumpText(bin []uint32) []ins.Instruction {
	result := make([]ins.Instruction, 0, len(bin))
	for _, bits := range bin {
		result = append(result, ins.Parse(bits))
	}
	return result
}
