package interpreter

func (vm *VirtualMachine) _CLS() {
	vm.Display.Clear()
}

func (vm *VirtualMachine) _RET() {
	// TODO: implement
}

func (vm *VirtualMachine) _JP(addr uint16) {
	vm.pc = addr
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
	vm.Display.UpdateState(&vm.memory, vm.i, vx, vy, n)
	vm.Display.Render()
}