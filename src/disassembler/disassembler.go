package main

import (
	"flag"
	"fmt"

	"encoding/binary"

	"github.com/abhinand20/emugo/common"
)

var InputFile string
const (
	instructionBytes = 2
	startAddr = 0x200
)

type Instruction struct {
	Opcode uint16
	Address uint16
	Name string
	LeftOp string
	RightOp string
}

func (inst *Instruction) print() {
	fmtStr := fmt.Sprintf("%04X: %04x %s", inst.Address, inst.Opcode, inst.Name)
	if len(inst.LeftOp) > 0 {
		fmtStr = fmt.Sprintf("%s  %s", fmtStr, inst.LeftOp)
		if len(inst.RightOp) > 0 {
			fmtStr = fmt.Sprintf("%s,%s", fmtStr, inst.RightOp)
		}
	}
	fmt.Printf("%s\n", fmtStr)
}

func initFlags() {
	flag.StringVar(&InputFile, "file", "", "File containing CHIP-8 hex code.")
}

func validateFlags() error {
	if len(InputFile) == 0 {
		return fmt.Errorf("input file not provided")
	}
	return nil
}

func parseHexInstruction(inst []byte, idx int) Instruction {
	var instruction Instruction
	upperNibble := (inst[0] >> 4)
	instruction.Opcode = binary.BigEndian.Uint16(inst)
	instruction.Address = startAddr + uint16(idx)
	switch upperNibble {
	case 0x00: {
		lowerByte := inst[1]
		switch lowerByte {
		case 0xE0: {
			instruction.Name = "CLS"
		}
		case 0xEE: {
			instruction.Name = "RET"	
		}
		default: instruction.Name = "UNK 0"
		}
	}
	case 0x01: {
		instruction.Name = "JP"
		leftOp := instruction.Opcode & (0xFFF)
		instruction.LeftOp = fmt.Sprintf("%X", leftOp)
	}
	case 0x06: {
		instruction.Name = "LD"
		leftOp := inst[0] & 0x0F
		rightOp := instruction.Opcode & 0xFF
		instruction.LeftOp = fmt.Sprintf("V%d", leftOp)
		instruction.RightOp = fmt.Sprintf("%X", rightOp)
	}
	case 0x07: {
		instruction.Name = "ADD"
		leftOp := inst[0] & 0x0F
		rightOp := instruction.Opcode & 0xFF
		instruction.LeftOp = fmt.Sprintf("V%d", leftOp)
		instruction.RightOp = fmt.Sprintf("%X", rightOp)	
	}
	case 0x0A: {
		instruction.Name = "LD"
		rightOp := instruction.Opcode & 0xFFF
		instruction.LeftOp = "I"
		instruction.RightOp = fmt.Sprintf("%X", rightOp)
	}
	case 0x0D: {
		instruction.Name = "DRW"
		x := inst[0] & 0x0F
		y := (inst[1] >> 4)
		n := inst[1] & 0x0F
		instruction.LeftOp = fmt.Sprintf("V%d", x)
		instruction.RightOp = fmt.Sprintf("V%d,%d", y, n)
	}
	default: instruction.Name = "UNK"
	}
	return instruction
}

func parseHexInstructions(arr []byte) []Instruction {
	var instructions []Instruction
	idx := 0
	end := len(arr)
	if end % 2 != 0 {
		panic("Something is wrong with the CHIP8 program.")
	}
	for idx < end {
		opcode := make([]byte, 2)
		opcode[0] = arr[idx]
		opcode[1] = arr[idx+1]
		inst := parseHexInstruction(opcode, idx)
		instructions = append(instructions, inst)
		idx += 2
	}
	return instructions
}

func main() {
	initFlags()
	flag.Parse()
	if err := validateFlags(); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	content, err := common.ReadFile(InputFile)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	instructs := parseHexInstructions(content)
	for _, i := range instructs {
		i.print()
	}
}