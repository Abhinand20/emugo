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
	if vm.sp != 0 {
		vm.sp += 1
	}
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
	if collision {
		vm.setVF()
		return
	}
	vm.resetVF()
}