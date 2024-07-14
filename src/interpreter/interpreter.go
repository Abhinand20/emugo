package interpreter

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"

	common "github.com/abhinand20/emugo/common"
	disp "github.com/abhinand20/emugo/display"
	"github.com/abhinand20/emugo/input"
)

var spriteData = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// VirtualMachine loads up the program into memory
// and executes it
type VirtualMachine struct {
	memory [4096]byte
	Display disp.Display
	pc uint16
	i uint16
	dt uint8
	ds uint8
	r [16]uint8
	Clk *time.Ticker
	sp uint16
	stack [16]uint16
	Keyboard *input.Keyboard
	/* States useful for debug mode */
	Debug bool
}

func (vm *VirtualMachine) Init(program []byte, clkSpeed int) {
	for idx := range program {
		vm.memory[common.ProgramStoreOffsetBytes + idx] = program[idx]
	}
	vm.loadSpritesInMemory()
	vm.pc = common.ProgramStoreOffsetBytes
	vm.Display.Init()
	vm.Clk = time.NewTicker(time.Second / time.Duration(clkSpeed))
}

func (vm *VirtualMachine) loadSpritesInMemory() {
	for idx := range spriteData {
		vm.memory[common.SpriteStartOffsetBytes + idx] = spriteData[idx]
	}
}

// Run is the main entry point for the VM
// it repeatedly goes through the fetch/execute cycle
func (vm *VirtualMachine) Run() error {
	vm.Keyboard.Start()
	for {
		// Wait for tick before proceeding
		<- vm.Clk.C
		if vm.Debug {
			instrBytes := vm.memory[vm.pc : vm.pc+2]
			debugInst := common.ParseHexInstruction(instrBytes, int(vm.pc))
			fmt.Print("> ")
			debugInst.Print()
			reader := bufio.NewReader(os.Stdin)
			// TODO(abhinandj): Add support for breakpoints.
			for {
				input, _ := reader.ReadString('\n')
				if strings.ToLower(input) == "n\n" {
					break
				}
			}
		}
		instruction, end := vm.fetch()
		if end {
			vm.Clk.Stop()
			break
		}
		err := vm.execute(instruction)
		if err != nil {
			return fmt.Errorf("could not execute instruction: %v", err)
		}
	}
	vm.Keyboard.Stop()
	return nil
}


func (vm *VirtualMachine) fetch() (*common.Opcode, bool) {
	if vm.pc >= uint16(len(vm.memory)) {
		return nil, true
	}
	instrBytes := vm.memory[vm.pc : vm.pc+2]
	opcode := binary.BigEndian.Uint16(instrBytes)
	vm.pc += 2
	return common.ParseOpcode(opcode), false
}

func (vm *VirtualMachine) execute(opcode *common.Opcode) error {
	switch opcode.NibbleUpper {
	case 0x00: {
		switch opcode.LowerByte {
		case 0xE0: vm._CLS()
		case 0xEE: vm._RET()
		default: return common.UnknownOpcodeErr(opcode.Opcode)
		}
	}
	case 0x01: vm._JP(opcode.Addr)
	case 0x02: vm._CALL(opcode.Addr)
	case 0x03: vm._SEVal(opcode.NibbleX, opcode.LowerByte)
	case 0x04: vm._SNEVal(opcode.NibbleX, opcode.LowerByte)
	case 0x05: vm._SE(opcode.NibbleX, opcode.NibbleY)
	case 0x09: vm._SNE(opcode.NibbleX, opcode.NibbleY)
	case 0x06: vm._LDVal(opcode.NibbleX, opcode.LowerByte)
	case 0x07: vm._ADDVal(opcode.NibbleX, opcode.LowerByte)
	case 0x08: {
		switch opcode.NibbleLower {
			case 0x00: vm._LD(opcode.NibbleX, opcode.NibbleY)
			case 0x01: vm._OR(opcode.NibbleX, opcode.NibbleY)
			case 0x02: vm._AND(opcode.NibbleX, opcode.NibbleY)
			case 0x03: vm._XOR(opcode.NibbleX, opcode.NibbleY)
			case 0x04: vm._ADD(opcode.NibbleX, opcode.NibbleY)
			case 0x05: vm._SUB(opcode.NibbleX, opcode.NibbleY)
			case 0x06: vm._SHR(opcode.NibbleX)
			case 0x07: vm._SUBN(opcode.NibbleX, opcode.NibbleY)
			case 0x0E: vm._SHL(opcode.NibbleX)
			default: return common.UnknownOpcodeErr(opcode.Opcode)
		}
	}
	case 0x0A: vm._LDI(opcode.Addr)
	case 0x0D: vm._DRW(opcode.NibbleX, opcode.NibbleY, opcode.NibbleLower)
	case 0x0F: {
		switch opcode.LowerByte {
			case 0x1E: vm._ADDI(opcode.NibbleX)
			case 0x33: vm._LDBCD(opcode.NibbleX)
			case 0x55: vm._STR(opcode.NibbleX)
			case 0x65: vm._LDR(opcode.NibbleX)
			default: return common.UnknownOpcodeErr(opcode.Opcode)
		}
	}
	default: return common.UnknownOpcodeErr(opcode.Opcode)
	}
	return nil
}

func (vm *VirtualMachine) setVF() {
	vm.r[0xF] = 1
}

func (vm *VirtualMachine) resetVF() {
	vm.r[0xF] = 0
}

func (vm *VirtualMachine) isVFSet() bool {
	return vm.r[0xF] == 1
}

func (vm *VirtualMachine) isOverflow(x, y byte) bool {
	res := x + y
	return !((res > x) == (y > 0))
}