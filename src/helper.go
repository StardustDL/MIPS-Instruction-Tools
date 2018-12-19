package main

import (
    "fmt"
    "bufio"
    "io"
    "os"
    "strings"
    "strconv"

    ins "./instruction"
)


func toASMs(instrs []ins.Instruction) []string {
	result := make([]string, len(instrs))
	for i, instr := range instrs {
		result[i] = instr.ToASM()
	}
	return result
}

func toBitStrings(bin []uint32) []string {
	result := make([]string, len(bin))
	for i, bits := range bin {
		result[i] = fmt.Sprintf("%08s", strconv.FormatUint(uint64(bits), 16))
	}
	return result
}

func toMIF(bin []uint8) []string {
	result := make([]string, 0, len(bin))
	result = append(result, "WIDTH=8;")
	result = append(result, fmt.Sprintf("DEPTH=%d;", len(bin)))
	result = append(result, "ADDRESS_RADIX=HEX;")
	result = append(result, "DATA_RADIX=HEX;")
	result = append(result, "CONTENT BEGIN")
	for i, bits := range bin {
		result = append(result, fmt.Sprintf("    %04X : %04X;", i, bits))
	}
	result = append(result, "END;")
	return result
}

func readAllLines(path string) ([]string, error) {
    file, err := os.OpenFile(path, os.O_RDONLY, 0666)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    _, err = file.Stat()
    if err != nil {
        return nil, err
    }

    buf := bufio.NewReader(file)

    result := make([]string, 0)

    for {
        line, _, err := buf.ReadLine()
        str := strings.TrimSpace(string(line))
        if err != nil {
            if err == io.EOF {
                break
            } else {
                return result, err
            }
        }
        result = append(result, str)
    }
    return result, nil
}

func writeAllLines(path string, content []string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, str := range content {
        _, err := file.WriteString(str + "\n")
        if err != nil {
            return err
        }
    }
    file.Sync()
    return nil
}

func readAllBytes(path string) ([]byte, error) {
    file, err := os.OpenFile(path, os.O_RDONLY, 0666)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    _, err = file.Stat()
    if err != nil {
        return nil, err
    }

    buf := bufio.NewReader(file)

    result := make([]byte, 0)

    for {
        bt, err := buf.ReadByte()
        if err != nil {
            if err == io.EOF {
                break
            } else {
                return result, err
            }
        }
        result = append(result, bt)
    }
    return result, nil
}

func writeAllBytes(path string, data []byte) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(data)
    return err
}
