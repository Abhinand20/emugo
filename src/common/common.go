package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	ProgramDir = "../roms/"
	ProgramReadOffsetBytes = 0
	ProgramStoreOffsetBytes = 512
	SpriteStartOffsetBytes = 0
)

type Opcode struct {
	Opcode uint16
	NibbleLower byte
	NibbleY byte 
	NibbleX byte
	NibbleUpper byte 
	LowerByte byte
	Addr uint16
}

func (o Opcode) String() {
	fmt.Printf("Opcode: %X\n", o.Opcode)
}


func ParseOpcode(opcode uint16) *Opcode {
	nibbleLower := extractNibble(opcode, 0)
	nibbleY := extractNibble(opcode, 1)
	nibbleX := extractNibble(opcode, 2)
	nibbleUpper := extractNibble(opcode, 3)
	lowerByte := extractLowerByte(opcode)
	addr := extractAddress(opcode)
	return &Opcode{
		Opcode: opcode,
		NibbleLower: nibbleLower,
		NibbleY: nibbleY,
		NibbleX: nibbleX,
		NibbleUpper: nibbleUpper,
		LowerByte: lowerByte,
		Addr: addr,
	}
}

func extractNibble(opcode uint16, idx byte) byte {
	var mask uint16 = 0xF
	return byte(opcode >> (idx * 4) & mask)
}

func extractLowerByte(opcode uint16) byte {
	var mask uint16 = 0xFF
	return byte(opcode & mask)
}

func extractAddress(opcode uint16) uint16 {
	var mask uint16 = 0xFFF
	return opcode & mask
}

func UnknownOpcodeErr(opcode uint16) error {
	return fmt.Errorf("unknown opcode: %X", opcode)
}


func ReadFile(file string) ([]byte, error) {
	absPath, err := filepath.Abs(ProgramDir + file)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file path '%s': %v", absPath, err)
	}
	f, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file '%s': %v", file, err)
	}
	defer f.Close()
	fs, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to stat file '%s': %v", file, err)
	}
	_, err = f.Seek(ProgramReadOffsetBytes, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("unable to find program start '%s': %v", file, err)
	}
	content := make([]byte, fs.Size() - ProgramReadOffsetBytes)
	n, err := f.Read(content)
	if err != nil {
		return content, fmt.Errorf("unable to read file '%s': %v", file, err)
	}
	fmt.Printf("Read %d bytes\n", n)
	return content, nil
}

func PrintHex(arr []byte) {
	fmt.Println("Num bytes = ", len(arr))
	for _, a := range arr {
		fmt.Printf("%x ", a)
	}
	fmt.Printf("\n")
}
