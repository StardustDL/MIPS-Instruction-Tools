package exec

import (
	"../cpu"
)

func add(it rinstr) {
	cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rs)+cpu.GetGPR(it.Rt))
}

func addu(it rinstr) {
	add(it)
}

func addi(it iinstr) {
	cpu.SetGPR(it.Rt, cpu.GetGPR(it.Rs)+signext16(it.Imm))
}

func addiu(it iinstr) {
	addi(it)
}

func sub(it rinstr) {
	cpu.SetGPR(it.Rd, cpu.GetGPR(it.Rs)-cpu.GetGPR(it.Rt))
}

func subu(it rinstr) {
	sub(it)
}

func lui(it iinstr) {
	cpu.SetGPR(it.Rt, uint32(it.Imm)<<16)
}

func mult(it rinstr) {
	vd := int64(signext64(cpu.GetGPR(it.Rs)))
	vr := int64(signext64(cpu.GetGPR(it.Rt)))
	cpu.SetAcc(uint64(vd * vr))
}

func multu(it rinstr) {
	vd := uint64(cpu.GetGPR(it.Rs))
	vr := uint64(cpu.GetGPR(it.Rt))
	cpu.SetAcc(vd * vr)
}

func div(it rinstr) {
	vd := int64(signext64(cpu.GetGPR(it.Rs)))
	vr := int64(signext64(cpu.GetGPR(it.Rt)))
	cpu.SetLO(uint32(vd / vr))
	cpu.SetHI(uint32(vd % vr))
}

func divu(it rinstr) {
	vd := uint64(cpu.GetGPR(it.Rs))
	vr := uint64(cpu.GetGPR(it.Rt))
	cpu.SetLO(uint32(vd / vr))
	cpu.SetHI(uint32(vd % vr))
}

func mfhi(it rinstr) {
	cpu.SetGPR(it.Rd, cpu.GetHI())
}

func mflo(it rinstr) {
	cpu.SetGPR(it.Rd, cpu.GetLO())
}

func mthi(it rinstr) {
	cpu.SetHI(cpu.GetGPR(it.Rs))
}

func mtlo(it rinstr) {
	cpu.SetLO(cpu.GetGPR(it.Rs))
}
