package emulator

const MEMORY_SIZE uint32 = 0x2000

var memory [MEMORY_SIZE]uint8

var _MASK_BYTE = [5]uint32{0x0, 0xff, 0xffff, 0xffffff, 0xffffffff}

func memoryRead(addr uint32, len uint8) (uint32, bool) {
	if !(len == 1 || len == 2 || len == 4) {
		return 0, false
	}
	if addr+uint32(len) >= MEMORY_SIZE {
		return 0, false
	}
	var result uint32 = 0
	for i := uint8(0); i < len; i++ {
		result |= uint32(memory[addr+uint32(i)]) << (uint32(i) << 3)
	}
	return result & _MASK_BYTE[len], true
}

func memoryWrite(addr uint32, len uint8, val uint32) bool {
	if !(len == 1 || len == 2 || len == 4) {
		return false
	}
	if addr+uint32(len) > MEMORY_SIZE {
		return false
	}
	val &= _MASK_BYTE[len]
	for i := uint8(0); i < len; i++ {
		memory[addr+uint32(i)] = uint8(val >> (uint32(i) << 3) & _MASK_BYTE[1])
	}
	return true
}
