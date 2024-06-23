package display

// A display interface that can be implemented
// using any display library under the hood
type Display interface {
	// Initial setup for the display
	Init()
	Clear()
	// Update the internal states to setup rendering
	// based on the draw instruction operands
	UpdateState(*[4096]byte, uint16, byte, byte, byte)
	// Reneder called for each display instruction execution
	Render()
}