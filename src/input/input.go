package input

import (
	"strings"
	"sync"

	"github.com/eiannone/keyboard"
)


type Keyboard struct {
	currentKey uint8
	isKeyHeld bool
	mu sync.Mutex
}

var keyMap = map[string]uint8{
	"0": 0x0,
	"1": 0x1,
	"2": 0x2,
	"3": 0x3,
	"4": 0x4,
	"5": 0x5,
	"6": 0x6,
	"7": 0x7,
	"8": 0x8,
	"9": 0x9,
	"a": 0x10,
	"b": 0x11,
	"c": 0x12,
	"d": 0x13,
	"e": 0x14,
	"f": 0x15,
}

func (kb *Keyboard) Start() {
	keyboard.Open()
	kb.currentKey = 255
	go kb.listner()
}

func (kb *Keyboard) Stop() {
	keyboard.Close()
}

// Returns the current pressed key, and whether it is being held down.
func (kb *Keyboard) GetPressedKey() (uint8, bool) {
	return kb.currentKey, kb.isKeyHeld
}

func (kb *Keyboard) listner() {
	for {
		charRune, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
        if key == keyboard.KeyCtrlC {
			keyboard.Close()
			break
		}
		pressedChar := strings.ToLower(string(charRune))
		if charIdx, ok := keyMap[pressedChar]; ok {
			kb.mu.Lock()
			if charIdx == kb.currentKey {
				kb.isKeyHeld = true
			} else {
				kb.isKeyHeld = false
			}
			kb.currentKey = charIdx
			kb.mu.Unlock()
		}
	}
}