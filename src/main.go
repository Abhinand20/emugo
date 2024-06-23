package main

import (
	"flag"
	"fmt"

	common "github.com/abhinand20/emugo/common"
	disp "github.com/abhinand20/emugo/display"
	"github.com/abhinand20/emugo/interpreter"
)

var InputFile string


func initFlags() {
	flag.StringVar(&InputFile, "file", "", "File containing CHIP-8 hex code.")
}

func validateFlags() error {
	if len(InputFile) == 0 {
		return fmt.Errorf("input file not provided")
	}
	return nil
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
	d := &disp.TerminalDisplay{
		Height: 32,
		Width: 64,
	}
	d.Init()
	vm := interpreter.VirtualMachine{
		Display: d,
	}
	vm.Init(content)
	if err := vm.Run(); err != nil {
		panic(err)
	}

}