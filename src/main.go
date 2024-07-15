package main

import (
	"flag"
	"fmt"

	common "github.com/abhinand20/emugo/common"
	disp "github.com/abhinand20/emugo/display"
	"github.com/abhinand20/emugo/input"
	"github.com/abhinand20/emugo/interpreter"
)

var inputFile string
var clkSpeed int
var debug bool

func initFlags() {
	flag.StringVar(&inputFile, "file", "", "File containing CHIP-8 hex code.")
	flag.IntVar(&clkSpeed, "clock_speed", 700, "Clock speed of the emulator in Hz.")
	flag.BoolVar(&debug, "debug", false, "Run debugger.")
}

func validateFlags() error {
	if len(inputFile) == 0 {
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
	content, err := common.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	d := &disp.TerminalDisplay{
		Height: 32,
		Width: 64,
	}
	d.Init()
	kb := &input.Keyboard{}
	vm := interpreter.VirtualMachine{
		Display: d,
		Keyboard: kb,
		Debug: debug,
	}
	vm.Init(content, clkSpeed)
	if debug {
		fmt.Println("Running debugger...\nEnter 'n' to step through instructions!")
	}
	// kb.Start()
	// time.Sleep(time.Second * 10)
	// kb.Stop()
	if err := vm.Run(); err != nil {
		panic(err)
	}
}