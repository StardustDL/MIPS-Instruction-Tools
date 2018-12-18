package main

import (
    "bufio"
    "io"
    "os"
    "strings"
)

func readAllLines(path string) ([]string, error) {
    file, err := os.OpenFile(path, os.O_RDONLY, 0666)
    if err != nil {
        return make([]string, 0), err
    }
    defer file.Close()

    _, err = file.Stat()
    if err != nil {
        return make([]string, 0), err
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
        return make([]byte, 0), err
    }
    defer file.Close()

    _, err = file.Stat()
    if err != nil {
        return make([]byte, 0), err
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
