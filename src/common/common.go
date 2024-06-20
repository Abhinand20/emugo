package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	ProgramDir = "../../roms/"
	ProgramOffsetBytes = 0 // 0x200 (512) bits
)

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
	_, err = f.Seek(ProgramOffsetBytes, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("unable to find program start '%s': %v", file, err)
	}
	content := make([]byte, fs.Size() - ProgramOffsetBytes)
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
