package exec

import (
    "../cpu"
    "../memory"
)

func lb(it iinstr) {
    cpu.SetGPR(it.Rt, signext(memory.Read(cpu.GetGPR(it.Rs)+signext16(it.Imm), 1), 2))
}

func lbu(it iinstr) {
    cpu.SetGPR(it.Rt, memory.Read(cpu.GetGPR(it.Rs)+signext16(it.Imm), 1))
}

func lh(it iinstr) {
    cpu.SetGPR(it.Rt, signext(memory.Read(cpu.GetGPR(it.Rs)+signext16(it.Imm), 2), 2))
}

func lhu(it iinstr) {
    cpu.SetGPR(it.Rt, memory.Read(cpu.GetGPR(it.Rs)+signext16(it.Imm), 2))
}

func lw(it iinstr) {
    cpu.SetGPR(it.Rt, memory.Read(cpu.GetGPR(it.Rs)+signext16(it.Imm), 4))
}

func sb(it iinstr) {
    memory.Write(cpu.GetGPR(it.Rs)+signext16(it.Imm), 1, cpu.GetGPR(it.Rt))
}

func sh(it iinstr) {
    memory.Write(cpu.GetGPR(it.Rs)+signext16(it.Imm), 2, cpu.GetGPR(it.Rt))
}

func sw(it iinstr) {
    memory.Write(cpu.GetGPR(it.Rs)+signext16(it.Imm), 4, cpu.GetGPR(it.Rt))
}