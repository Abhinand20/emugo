package interpreter

import (
	"encoding/binary"
	"fmt"

	common "github.com/abhinand20/emugo/common"
	disp "github.com/abhinand20/emugo/display"
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
	// TODO: Add support for a stack/SP
}

func (vm *VirtualMachine) Init(program []byte) {
	for idx := range program {
		vm.memory[common.ProgramStoreOffsetBytes + idx] = program[idx]
	}
	vm.loadSpritesInMemory()
	vm.pc = common.ProgramStoreOffsetBytes
	vm.Display.Init()
}

func (vm *VirtualMachine) loadSpritesInMemory() {
	for idx := range spriteData {
		vm.memory[common.SpriteStartOffsetBytes + idx] = spriteData[idx]
	}
}

// Run is the main entry point for the VM
// it repeatedly goes through the fetch/execute cycle
func (vm *VirtualMachine) Run() error {
	for {
		// TODO: Adjust once frequency is implemented
		instruction, done := vm.fetch()
		if done {
			fmt.Println("Done!")
			break
		}
		err := vm.execute(instruction)
		if err != nil {
			return fmt.Errorf("could not execute instruction: %v", err)
		}
	}
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
	case 0x06: vm._LDVal(opcode.NibbleX, opcode.LowerByte)
	case 0x07: vm._ADDVal(opcode.NibbleX, opcode.LowerByte)
	case 0x0A: vm._LDI(opcode.Addr)
	case 0x0D: vm._DRW(opcode.NibbleX, opcode.NibbleY, opcode.NibbleLower)
	default: return common.UnknownOpcodeErr(opcode.Opcode)
	}
	return nil
}

func (vm *VirtualMachine) setVF() {
	vm.r[len(vm.r) - 1] = 1
}

func (vm *VirtualMachine) resetVF() {
	vm.r[len(vm.r) - 1] = 0
}

func (vm *VirtualMachine) isVFSet() bool {
	return vm.r[len(vm.r) - 1] == 1
}