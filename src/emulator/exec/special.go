package exec

import (
	"../../instruction"
	"../cpu"
)

func nop(it rinstr) {

}

func syscall(it rinstr) {
	doSystemCall(cpu.GetGPR(instruction.GPR_V0))
}
