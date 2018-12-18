package exec

import (
    "../cpu"
)

func sll(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rt)<<it.Shamt)
}

func sllv(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rt)<<cpu.GetGPR(it.Rs))
}

func sra(it rinstr) {
    cpu.SetGPR(it.Rd, uint32(int32(cpu.GetGPR(it.Rt))>>uint32(it.Shamt)))
}

func srav(it rinstr) {
    cpu.SetGPR(it.Rd, uint32(int32(cpu.GetGPR(it.Rt))>>cpu.GetGPR(it.Rs)))
}

func srl(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rt)>>uint32(it.Shamt))
}

func srlv(it rinstr) {
    cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rt)>>cpu.GetGPR(it.Rs))
}
