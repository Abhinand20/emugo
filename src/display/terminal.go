package display

import (
	"fmt"
)

// Simple terminal display implements the Display interface
type TerminalDisplay struct {
	Grid [][]uint8
	Height uint32
	Width uint32
}

func (t *TerminalDisplay) Init() {
	t.Grid = make([][]uint8, t.Height)
	for i := range t.Grid {
		t.Grid[i] = make([]uint8, t.Width)
	}
	t.resetGrid()
}

func (t *TerminalDisplay) Clear() {
	t.resetGrid()
}

func (t *TerminalDisplay) UpdateState(memory *[4096]byte, i uint16, vx, vy, n byte) bool {
	// TODO: Handle 1) wrapping behavior 2) Setting VF for collisions
	var ib uint8 = 0
	var collision bool = false
	for ib < n {
		spriteRow := memory[i + uint16(ib)]
		// XOR with all pixels in this row
		var idx uint8 = 0
		for idx < 8 {
			currRow := vy + ib
			currCol := vx + idx
			prevSet := t.Grid[currRow][currCol]
			t.Grid[currRow][currCol] ^= ((spriteRow >> (7 - idx)) & 0x1)
			if prevSet == 1 && t.Grid[currRow][currCol] == 0 {
				collision = true
			}
			idx++
		}
		ib++
	}
	return collision
}

func (t *TerminalDisplay) Render() {
	// Hacky way to clear terminal in macOS/linux, won't work on windows.
	fmt.Print("\033[H\033[2J")
	t.drawHorizontalBorder()
	for i := range t.Grid {
		fmt.Print("|")
		for j := range t.Grid[i] {
			displayChar := ' '
			if t.Grid[i][j] == 1 {
				displayChar = '*'
			}
			fmt.Printf("%c", displayChar)
		}
		fmt.Printf("|\n")
	}
	t.drawHorizontalBorder()
}

func (t *TerminalDisplay) drawHorizontalBorder() {
	fmt.Print("+")
	for r := 0; r < int(t.Width); r++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func (t *TerminalDisplay) resetGrid() {
	for i := range t.Grid {
		for j := range t.Grid[i] {
			t.Grid[i][j] = 0
		}
	}
}