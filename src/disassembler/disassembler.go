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
)


func initFlags() {
	flag.StringVar(&InputFile, "file", "", "File containing CHIP-8 hex code.")
}

func validateFlags() error {
	if len(InputFile) == 0 {
		return fmt.Errorf("input file not provided")
	}
	return nil
}

func parseHexInstructions(arr []byte) []common.Instruction {
	var instructions []common.Instruction
	idx := 0
	end := len(arr)
	for idx < end {
		var inst common.Instruction
		// Not a valid instruction
		if idx == end - 1 {
			inst.Name = "UNK"
			inst.Address = common.StartAddr + uint16(idx)
			inst.Opcode = binary.BigEndian.Uint16([]byte{arr[idx], 0})
		} else {
			opcode := make([]byte, 2)
			opcode[0] = arr[idx]
			opcode[1] = arr[idx+1]
			inst = common.ParseHexInstruction(opcode, idx)
		}
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
		i.Print()
	}
}