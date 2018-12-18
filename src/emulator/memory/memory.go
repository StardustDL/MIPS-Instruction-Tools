package memory

import (
    "fmt"
)

const MEMORY_SIZE uint32 = 0x2000

var memory [MEMORY_SIZE]uint8

var _MASK_BYTE = [5]uint32{0x0, 0xff, 0xffff, 0xffffff, 0xffffffff}

func Read(addr uint32, len uint8) uint32 {
    if !(len == 1 || len == 2 || len == 4) {
        panic(fmt.Sprintf("Memory rw with unexpected len %d", len))
    }
    if addr+uint32(len) > MEMORY_SIZE {
        panic(fmt.Sprintf("Memory rw with too large address %d", addr+uint32(len)))
    }
    var result uint32 = 0
    for i := uint8(0); i < len; i++ {
        result |= uint32(memory[addr+uint32(i)]) << (uint32(i) << 3)
    }
    return result & _MASK_BYTE[len]
}

func Write(addr uint32, len uint8, val uint32) {
    if !(len == 1 || len == 2 || len == 4) {
        panic(fmt.Sprintf("Memory rw with unexpected len %d", len))
    }
    if addr+uint32(len) > MEMORY_SIZE {
        panic(fmt.Sprintf("Memory rw with too large address %d", addr+uint32(len)))
    }
    val &= _MASK_BYTE[len]
    for i := uint8(0); i < len; i++ {
        memory[addr+uint32(i)] = uint8(val >> (uint32(i) << 3) & _MASK_BYTE[1])
    }
}
