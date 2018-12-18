package exec

import (
    "fmt"
    "bufio"
    "os"
    "strconv"

    "../../instruction"
    "../cpu"
    "../memory"
)

const (
    SYS_PRINT_INT    = uint32(1)
    SYS_PRINT_STRING = uint32(4)
    SYS_READ_INT     = uint32(5)
    SYS_READ_STRING  = uint32(8)
    SYS_EXIT         = uint32(10)
    SYS_PRINT_CHAR   = uint32(11)
)

func readString(addr uint32) string {
    bytes := make([]byte, 0)
    chr := memory.Read(addr, 1)
    for addr++; chr != 0; addr++ {
        bytes = append(bytes, byte(chr))
        chr = memory.Read(addr, 1)
    }
    return string(bytes)
}

func writeString(addr uint32, bufsize uint32, str string) {
    if uint32(len(str)) > bufsize-1 {
        str = str[0 : bufsize-1]
    }
    for i, chr := range str {
        memory.Write(addr+uint32(i), 1, uint32(chr))
    }
    memory.Write(addr+uint32(len(str)), 1, uint32(0))
}

var stdin *bufio.Scanner

func doSystemCall(id uint32) {
    switch id {
    case SYS_EXIT:
        if IsDebug{
            println("Exited")
        }
        State = MEMU_EXITED
    case SYS_PRINT_INT:
        val := cpu.GetGPR(instruction.GPR_A0)
        print(val)
    case SYS_PRINT_STRING:
        strpos := cpu.GetGPR(instruction.GPR_A0)
        str := readString(strpos)
        print(str)
    case SYS_PRINT_CHAR:
        val := cpu.GetGPR(instruction.GPR_A0)
        print(string(rune(val)))
    case SYS_READ_INT:
        if stdin == nil {
            stdin = bufio.NewScanner(os.Stdin)
        }
        if stdin.Scan() {
            str := stdin.Text()
            val, err := strconv.ParseInt(str, 0, 32)
            if err == nil {
                cpu.SetGPR(instruction.GPR_V0, uint32(val))
                break
            }
        }
        cpu.SetGPR(instruction.GPR_V0, uint32(0))
    case SYS_READ_STRING:
        strpos := cpu.GetGPR(instruction.GPR_A0)
        buflen := cpu.GetGPR(instruction.GPR_A1)
        if stdin == nil {
            stdin = bufio.NewScanner(os.Stdin)
        }
        if stdin.Scan() {
            str := stdin.Text()
            writeString(strpos, buflen, str)
        } else {
            writeString(strpos, buflen, "")
        }
    }
}

func doBreak(code uint32){
    if IsDebug{
        
    }

    println(fmt.Sprintf("Recieve break code: %d",code))
}