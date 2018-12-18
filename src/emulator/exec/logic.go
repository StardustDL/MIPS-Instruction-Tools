package exec

import (
    "../cpu"
)

func and(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rs)&cpu.GetGPR(it.Rt))
}

func andi(it iinstr) {
    cpu.SetGPR(it.Rt, cpu.GetGPR(it.Rs)&uint32(it.Imm))
}

func or(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rs)|cpu.GetGPR(it.Rt))
}

func ori(it iinstr) {
    cpu.SetGPR(it.Rt, cpu.GetGPR(it.Rs)|uint32(it.Imm))
}

func nor(it rinstr) {
    cpu.SetGPR(it.Rd, ^(cpu.GetGPR(it.Rs)|cpu.GetGPR(it.Rt)))
}

func xor(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rs)^cpu.GetGPR(it.Rt))
}

func xori(it iinstr) {
    cpu.SetGPR(it.Rt, cpu.GetGPR(it.Rs)^uint32(it.Imm))
}
