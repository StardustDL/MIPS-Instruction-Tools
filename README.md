# MIPS Instruction Tools

[![](https://img.shields.io/github/stars/StardustDL/MIPS-Instruction-Tools.svg?style=social&label=Stars)](https://github.com/StardustDL/MIPS-Instruction-Tools) [![](https://img.shields.io/github/forks/StardustDL/MIPS-Instruction-Tools.svg?style=social&label=Fork)](https://github.com/StardustDL/MIPS-Instruction-Tools) ![](http://progressed.io/bar/60?title=developing) [![](https://img.shields.io/github/license/StardustDL/MIPS-Instruction-Tools.svg)](https://github.com/StardustDL/MIPS-Instruction-Tools/blob/master/LICENSE)

An experimental tool for MIPS architecture (MIPS-32, Little-Endian). Now, this project contains two tools:

- Assembler for MIPS-32 architecture
- Simulator for MIPS-32 architecture
- Simple dumper for MIPS-32 architecture

Attention:

- This is an experimental project for learning purposes.
- To try these techniques: Go.
- Created for [Experimental-MIPS-CPU](https://github.com/StardustDL/Experimental-MIPS-CPU) 

Project Status:

|||
|-|-|
|Repository|[![](https://img.shields.io/github/issues/StardustDL/MIPS-Instruction-Tools.svg)](https://github.com/StardustDL/MIPS-Instruction-Tools/issues/) [![](https://img.shields.io/github/issues-pr/StardustDL/MIPS-Instruction-Tools.svg)](https://github.com/StardustDL/MIPS-Instruction-Tools/pulls/)|
|Dependencies|![](https://img.shields.io/librariesio/github/StardustDL/MIPS-Instruction-Tools.svg)|
|Build|[![](https://img.shields.io/travis/StardustDL/MIPS-Instruction-Tools/master.svg?label=master)](https://travis-ci.org/StardustDL/MIPS-Instruction-Tools) ![](https://img.shields.io/travis/StardustDL/MIPS-Instruction-Tools/dev.svg?label=dev)|
|Release|[![](https://img.shields.io/github/release-pre/StardustDL/MIPS-Instruction-Tools.svg)](https://github.com/StardustDL/MIPS-Instruction-Tools/releases/latest/) [![](https://img.shields.io/github/tag/StardustDL/MIPS-Instruction-Tools.svg)](https://github.com/StardustDL/MIPS-Instruction-Tools/tags)|

# Use

- Before use these commands, update `Makefile` to fit your system

```sh
make build
cd bin ; ./mif
```

# Instructions

See [here](https://www.mips.com/?do-download=mips32-instruction-set-quick-reference-v1-01) and [here](https://www.mips.com/products/architectures/mips32-3/)

Count: 51

- add
- addu
- addi
- addiu
- sub
- subu
- lui
- mul
- mult
- multu
- div
- divu
- mfhi
- mflo
- mthi
- mtlo
- slt
- sltu
- slti
- sltiu
- and
- andi
- or
- ori
- nor
- xor
- xori
- sll
- sllv
- sra
- srav
- srl
- srlv
- beq
- bne
- bgez
- blez
- bgtz
- bltz
- bgezal
- bltzal
- j
- jal
- jr
- jalr
- lb
- lbu
- lh
- lhu
- sb
- lw
- sw
- sh
- syscall
- break

## Assembler pseudo-instruction

Count: 14

- li
- la
- neg
- negu
- b
- bal
- beqz
- bnez
- nop
- move
- push
- pop
- call
- ret

## To append

None
