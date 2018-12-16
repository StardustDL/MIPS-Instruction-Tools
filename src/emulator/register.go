package emulator

var regs [32]uint32

var pc, npc uint32

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
