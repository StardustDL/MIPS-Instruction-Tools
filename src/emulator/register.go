package emulator

var regs [32]uint32

var pc, npc uint32

const (
	GPR_ZERO = uint8(iota)
	GPR_AT
	GPR_V0
	GPR_V1
	GPR_A0
	GPR_A1
	GPR_A2
	GPR_A3
	GPR_T0
	GPR_T1
	GPR_T2
	GPR_T3
	GPR_T4
	GPR_T5
	GPR_T6
	GPR_T7
	GPR_S0
	GPR_S1
	GPR_S2
	GPR_S3
	GPR_S4
	GPR_S5
	GPR_S6
	GPR_S7
	GPR_T8
	GPR_T9
	GPR_K0
	GPR_K1
	GPR_GP
	GPR_SP
	GPR_FP
	GPR_RA
)

func setRegister(id uint8, val uint32) bool {
	if !(0 <= id && id < 32) {
		return false
	}
	regs[id] = val
	return true
}

func getRegister(id uint8) (uint32, bool) {
	if !(0 <= id && id < 32) {
		return 0, false
	}
	return regs[id], true
}
