package emulator

import (
	"fmt"
	"runtime/debug"

	"../instruction"
	"./cpu"
	"./exec"
	"./memory"
)

func executeOne(instr instruction.Instruction) {

	defer handleErrorWhileExecuting()

	if exec.IsDebug {
		fmt.Printf("0x%x: %08x %s\n", cpu.PC, instr.ToBits(), instr.ToASM())
	}
	token := instr.GetToken()
	exc, ok := exec.ExecTable[token]
	if !ok {
		panic(fmt.Sprintf("No this instruction %s", token))
	}
	switch exc.(type) {
	case exec.ExecRFunc:
		ri, ok := instr.(instruction.RInstruction)
		if !ok {
			panic(fmt.Sprintf("The instruction type isn't fitting exec type"))
		}
		exc.(exec.ExecRFunc)(ri)
	case exec.ExecIFunc:
		ri, ok := instr.(instruction.IInstruction)
		if !ok {
			panic(fmt.Sprintf("The instruction type isn't fitting exec type"))
		}
		exc.(exec.ExecIFunc)(ri)
	case exec.ExecJFunc:
		ri, ok := instr.(instruction.JInstruction)
		if !ok {
			panic(fmt.Sprintf("The instruction type isn't fitting exec type"))
		}
		exc.(exec.ExecJFunc)(ri)
	default:
		panic(fmt.Sprintf("Internal error: exec type error"))
	}
}

func handleErrorWhileExecuting() {
	if err := recover(); err != nil {
		exec.State = exec.MEMU_ERROR
		fmt.Printf("Error %s\n", err.(string))
		debug.PrintStack()
	}
}

func Execute(entry uint32, isdebug bool) bool {
	defer handleErrorWhileExecuting()

	if exec.State != exec.MEMU_INITIALIZED {
		fmt.Println("Not initialized")
		return false
	}
	exec.IsDebug = isdebug
	exec.State = exec.MEMU_RUNNING

	exec.InitializePC(entry)

	for {
		bits := memory.Read(cpu.PC, 4)
		instr := instruction.Parse(bits)
		executeOne(instr)
		if exec.State != exec.MEMU_RUNNING {
			break
		}
		exec.UpdatePC()
	}

	return exec.State == exec.MEMU_EXITED
}

func Initialize(bin []uint8) bool {
	exec.InitializeTable()
	exec.State = exec.MEMU_EXITED
	for i := 0; i < 32; i++ {
		cpu.SetGPR(uint8(i), 0)
	}
	for i := uint32(0); i < memory.MEMORY_SIZE; i++ {
		memory.Write(uint32(i), 1, 0)
	}
	for i, bits := range bin {
		memory.Write(uint32(i), 1, uint32(bits))
	}
	exec.State = exec.MEMU_INITIALIZED
	return true
}

func ShowRegisters() {
	for i := 0; i < 32; i++ {
		val := cpu.GetGPR(uint8(i))
		fmt.Printf("$%02d %08x; ", i, val)
		if (i+1)%4 == 0 {
			println()
		}
	}
}
