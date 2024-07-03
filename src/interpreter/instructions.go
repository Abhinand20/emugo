package interpreter

func (vm *VirtualMachine) _CLS() {
	vm.Display.Clear()
}

// OPCODE: 00EE
func (vm *VirtualMachine) _RET() {
	vm.pc = vm.stack[vm.sp]
	if vm.sp != 0 {
		vm.sp -= 1
	}
}

// OPCODE: 2nnn
func (vm *VirtualMachine) _CALL(addr uint16) {
	if vm.sp >= uint16(len(vm.stack)) {
		panic("[Stack Overflow] Cannot execute instruction.")
	}
	vm.sp++
	vm.stack[vm.sp] = vm.pc
	vm.pc = addr
}

func (vm *VirtualMachine) _JP(addr uint16) {
	vm.pc = addr
}

// OPCODE: 3xkk
func (vm *VirtualMachine) _SEVal(x, kk byte) {
	if vm.r[x] == kk {
		vm.pc += 2
	}
}

// OPCODE: 4xkk
func (vm *VirtualMachine) _SNEVal(x, kk byte) {
	if vm.r[x] != kk {
		vm.pc += 2
	}
}

// OPCODE: 5xy0
func (vm *VirtualMachine) _SE(x, y byte) {
	if vm.r[x] == vm.r[y] {
		vm.pc += 2
	}
}


// OPCODE: 9xy0
func (vm *VirtualMachine) _SNE(x, y byte) {
	if vm.r[x] != vm.r[y] {
		vm.pc += 2
	}
}

func (vm *VirtualMachine) _LDVal(x, val byte) {
	vm.r[x] = val
}

// OPCODE: 8xy0
func (vm *VirtualMachine) _LD(x, y byte) {
	vm.r[x] = vm.r[y]
}

// OPCODE: 8xy1
func (vm *VirtualMachine) _OR(x, y byte) {
	vm.r[x] |= vm.r[y]
}

// OPCODE: 8xy2
func (vm *VirtualMachine) _AND(x, y byte) {
	vm.r[x] &= vm.r[y]
}

// OPCODE: 8xy3
func (vm *VirtualMachine) _XOR(x, y byte) {
	vm.r[x] ^= vm.r[y]
}

// OPCODE: 8xy4
func (vm *VirtualMachine) _ADD(x, y byte) {
	// Tricky edge case: if `VF` is passed as x
	vx := vm.r[x]
	vm.r[x] += vm.r[y]
	if vm.isOverflow(vx, vm.r[y]) {
		vm.setVF()
	} else {
		vm.resetVF()
	}
}

// OPCODE: 8xy5
func (vm *VirtualMachine) _SUB(x, y byte) {
	// Tricky edge case: if `VF` is passed as x
	vx := vm.r[x]
	vm.r[x] -= vm.r[y]
	if vx >= vm.r[y] {
		vm.setVF()
	} else {
		vm.resetVF()
	}
}

// OPCODE: 8xy6
func (vm *VirtualMachine) _SHR(x byte) {
	vx := vm.r[x]
	vm.r[x] >>= 1
	vm.resetVF()
	if (vx & 0x1) == 1 {
		vm.setVF()
	}
}

// OPCODE: 8xy7
func (vm *VirtualMachine) _SUBN(x, y byte) {
	// Tricky edge case: if `VF` is passed as x
	vx := vm.r[x]
	vm.r[x] = vm.r[y] - vm.r[x]
	if vm.r[y] >= vx {
		vm.setVF()
	} else {
		vm.resetVF()
	}
}

// OPCODE: 8xyE
func (vm *VirtualMachine) _SHL(x byte) {
	vx := vm.r[x]
	vm.r[x] <<= 1
	vm.resetVF()
	if (vx >> 7) & 0x1 == 1 {
		vm.setVF()
	}
}

// OPCODE: Fx65
func (vm *VirtualMachine) _LDR(x byte) {
	readAddr := vm.i
	for idx := byte(0); idx <= x; idx++ {
		vm.r[idx] = vm.memory[readAddr]
		readAddr++
	}
}

// OPCODE: Fx55
func (vm *VirtualMachine) _STR(x byte) {
	storeAddr := vm.i
	for idx := byte(0); idx <= x; idx++ {
		vm.memory[storeAddr] = vm.r[idx]
		storeAddr++
	}
}

// OPCODE: Fx33
func (vm *VirtualMachine) _LDBCD(x byte) {
	vm.memory[vm.i] = vm.r[x] / 100
	vm.memory[vm.i + 1] = (vm.r[x] / 10) % 10
	vm.memory[vm.i + 2] = vm.r[x] % 10
}


// OPCODE: Fx1E
func (vm *VirtualMachine) _ADDI(x byte) {
	vm.i += uint16(vm.r[x])
}


func (vm *VirtualMachine) _ADDVal(x, val byte) {
	vm.r[x] += val
}

func (vm *VirtualMachine) _LDI(addr uint16) {
	vm.i = addr
}

func (vm *VirtualMachine) _DRW(x, y, n byte) {
	vx := vm.r[x]
	vy := vm.r[y]
	collision := vm.Display.UpdateState(&vm.memory, vm.i, vx, vy, n)
	vm.Display.Render()
	vm.resetVF()
	if collision {
		vm.setVF()
	}
}