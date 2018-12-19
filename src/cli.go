package main

import (
	"flag"
	"fmt"
	"os"

	ass "./assembler"
	ins "./instruction"
	sim "./simulator"
)

func cliAs(inputFile string, bitsFile string, asmFile string, binFile string, mifFile string, dataSegment uint32, textSegment uint32, fullSize int32) (int, []ins.Instruction, *ass.AssembleResult) {
	if inputFile == "" {
		fmt.Printf("Please give the input file name")
		return -1, nil, nil
	}
	content, err := readAllLines(inputFile)
	if err != nil {
		fmt.Printf("File %s reading error: %v\n", inputFile, err)
		return -1, nil, nil
	}
	print("Assembling...")
	instrs, builded, err := ass.Assemble(content, ass.AssembleConfig{Data: dataSegment, Text: textSegment}, fullSize)
	if err == nil {
		println("done")
	} else {
		println("failed")
		println(err.Error())
		return -1, nil, nil
	}
	println("Instruction count:", len(instrs))
	fmt.Printf("Full segment: [0x%08x, 0x%08x), size: 0x%08x\n", builded.Full.Start, builded.Full.End, builded.Full.End-builded.Full.Start)
	fmt.Printf("Data segment: [0x%08x, 0x%08x), size: 0x%08x\n", builded.Data.Start, builded.Data.End, builded.Data.End-builded.Data.Start)
	fmt.Printf("Text segment: [0x%08x, 0x%08x), size: 0x%08x\n", builded.Text.Start, builded.Text.End, builded.Text.End-builded.Text.Start)
	if bitsFile != "" {
		err = writeAllLines(bitsFile, toBitStrings(ins.ToBin(instrs)))
		if err != nil {
			println("Generate bit string file failed", err)
			return -1, nil, nil
		}
		println("Bit string file:", bitsFile)
	}
	if asmFile != "" {
		err = writeAllLines(asmFile, toASMs(instrs))
		if err != nil {
			println("Generate asm file failed", err)
			return -1, nil, nil
		}
		println("ASM file:", asmFile)
	}
	if binFile != "" {
		err = writeAllBytes(binFile, builded.Bin)
		if err != nil {
			println("Generate asm file failed", err)
			return -1, nil, nil
		}
		println("Bin file:", binFile)
	}
	if mifFile != "" {
		err = writeAllLines(mifFile, toMIF(builded.Bin))
		if err != nil {
			println("Generate mif file failed", err)
			return -1, nil, nil
		}
		println("MIF file:", mifFile)
	}
	return 0, instrs, &builded
}

func usage() {
	fmt.Fprintf(os.Stderr, `mip version: mip/0.0.1
Usage: mip [options] as/sim/dump [inputFile]

Options:
`)
	flag.PrintDefaults()
}

func cliMain() int {
	var asmFile, binFile, bitsFile, mifFile, verb, inputFile string
	var textSegment, dataSegment uint64
	var fullSize, entry int64
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "Show help screen")
	flag.StringVar(&asmFile, "asm", "", "ASM file name")
	flag.StringVar(&binFile, "bin", "", "Bin file name")
	flag.StringVar(&mifFile, "mif", "", "Mif file name")
	flag.StringVar(&bitsFile, "bits", "", "Bit string file name")
	flag.Uint64Var(&textSegment, "text", 0, "Starting address of text segment")
	flag.Uint64Var(&dataSegment, "data", 0, "Starting address of data segment")
	flag.Int64Var(&entry, "entry", -1, "Program entry point, negtive for default (start of text segment)")
	flag.Int64Var(&fullSize, "size", -1, "Full size of program, negtive for no bin data")
	flag.Usage = usage

	flag.Parse()	

	if helpFlag {
		flag.Usage()
		return 0
	}

	verb = flag.Arg(0)
	inputFile = flag.Arg(1)

	switch verb {
	case "as":
		retcode, _, _ := cliAs(inputFile, bitsFile, asmFile, binFile, mifFile, uint32(dataSegment), uint32(textSegment), int32(fullSize))
		return retcode
	case "sim":
		var _entry uint32
		if binFile != "" {
			if entry < 0 {
				fmt.Printf("Must give entry point for bin file\n")
				return -1
			}
			content, err := readAllBytes(binFile)
			if err != nil {
				fmt.Printf("File %s reading error: %v\n", binFile, err)
				return -1
			}

			print("Initializing for simulating...")
			if !sim.Initialize(content, breakHandler, nil) {
				println("failed")
				return -1
			}

			_entry = uint32(entry)

			println("done")
		} else if asmFile != "" {
			retcode, _, buildedptr := cliAs(asmFile, "", "", "", "", uint32(dataSegment), uint32(textSegment), int32(fullSize))
			if retcode != 0 {
				return retcode
			}
			builded := *buildedptr
			if builded.Full.End == 0 {
				println("No bin data. Stop simulating")
				return -1
			}

			print("Initializing for simulating...")
			if !sim.Initialize(builded.Bin, breakHandler, nil) {
				println("failed")
				return -1
			}

			if entry < 0 {
				_entry = builded.Text.Start
			} else {
				_entry = uint32(entry)
			}

			println("done")
		} else {
			fmt.Printf("Please give the input file name.")
			return -1
		}
		println("Executing...")
		flg := sim.Execute(_entry, false)
		println("Executed:", flg)
		fmt.Println("Registers")
		sim.ShowRegisters()
	case "dump":
		fmt.Printf("Not support this function now")
		return -1
	default:
		println("No this verb.")
		return -1
	}
	return 0
}
