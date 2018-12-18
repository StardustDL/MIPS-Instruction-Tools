package exec

import (
    "../cpu"
)

func beq(it iinstr) {
    if cpu.GetGPR(it.Rs) == cpu.GetGPR(it.Rt) {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func bne(it iinstr) {
    if cpu.GetGPR(it.Rs) != cpu.GetGPR(it.Rt) {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func bgez(it iinstr) {
    if int32(cpu.GetGPR(it.Rs)) >= 0 {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func bgezal(it iinstr) {
    beforeCall()
    if int32(cpu.GetGPR(it.Rs)) >= 0 {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func bgtz(it iinstr) {
    if int32(cpu.GetGPR(it.Rs)) > 0 {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func blez(it iinstr) {
    if int32(cpu.GetGPR(it.Rs)) <= 0 {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func bltz(it iinstr) {
    if int32(cpu.GetGPR(it.Rs)) < 0 {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func bltzal(it iinstr) {
    beforeCall()
    if int32(cpu.GetGPR(it.Rs)) < 0 {
        jumpOneDelay(npc + (signext16(it.Imm) << 2))
    }
}

func j(it jinstr){
    jumpOneDelay((cpu.PC & 0xf0000000) | (it.Imm << 2))
}

func jal(it jinstr){
    beforeCall()
    j(it)
}

func jr(it rinstr){
    jumpOneDelay(cpu.GetGPR(it.Rs))
}

func jalr(it rinstr){
    tmp := cpu.GetGPR(it.Rs)
    beforeCallSet(it.Rd)
    jumpOneDelay(tmp)
}