package exec

import (
	"../cpu"
)

func slt(it rinstr) {
	if int32(cpu.GetGPR(it.Rs)) < int32(cpu.GetGPR(it.Rt)) {
		cpu.SetGPR(it.Rd, 1)
	} else {
		cpu.SetGPR(it.Rd, 0)
	}
}

func slti(it iinstr) {
	if int32(cpu.GetGPR(it.Rs)) < int32(signext16(it.Imm)) {
		cpu.SetGPR(it.Rt, 1)
	} else {
		cpu.SetGPR(it.Rt, 0)
	}
}

func sltu(it rinstr) {
	if cpu.GetGPR(it.Rs) < cpu.GetGPR(it.Rt) {
		cpu.SetGPR(it.Rd, 1)
	} else {
		cpu.SetGPR(it.Rd, 0)
	}
}

func sltiu(it iinstr) {
	if cpu.GetGPR(it.Rs) < uint32(it.Imm) {
		cpu.SetGPR(it.Rt, 1)
	} else {
		cpu.SetGPR(it.Rt, 0)
	}
}
