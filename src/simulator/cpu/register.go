package cpu

import (
    "fmt"
    "errors"
)

var regs [32]uint32

var PC uint32

var hi, lo uint32

func SetAcc(val uint64) {
    hi = uint32(val >> 32)
    lo = uint32(val & 0xffffffff)
}

func GetHI() uint32 {
    return hi
}

func GetLO() uint32 {
    return lo
}

func SetHI(val uint32){
    hi=val
}

func SetLO(val uint32){
    lo=val
}

func SetGPR(id uint8, val uint32) {
    if !(0 <= id && id < 32) {
        panic(errors.New(fmt.Sprintf("Register set failed %d", id)))
    }
    if id == 0 && val != 0{
        panic(errors.New(fmt.Sprintf("Try to set $zero to %d", val)))
    }
    regs[id] = val
}

func GetGPR(id uint8) uint32 {
    if !(0 <= id && id < 32) {
        panic(errors.New(fmt.Sprintf("Register get failed %d", id)))
    }
    return regs[id]
}
