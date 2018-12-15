package emulator

const MEMORY_SIZE = 1 << 10

var memory [MEMORY_SIZE]uint32

func memoryRead(addr uint32) (uint32, bool) {
	if addr >= MEMORY_SIZE {
		return 0, false
	}
	return memory[addr], true
}

func memoryWrite(addr uint32, val uint32) bool {
	if addr >= MEMORY_SIZE {
		return false
	}
	memory[addr] = val
	return true
}
