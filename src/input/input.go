package input

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)


type Keyboard struct {
	currentKeysPressed [16]bool
	prevKeysPressed [16]bool
	tempKeysPressed [16]bool
	// approximate key releases using delay between key press events
	keysDown map[byte]time.Time
	mu sync.Mutex
}

const (
	releaseDelay = time.Second / 5
	keyChannelSize = 20
)

var KeyMap = map[string]byte{
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
	"a": 0x0A,
	"b": 0x0B,
	"c": 0x0C,
	"d": 0x0D,
	"e": 0x0E,
	"f": 0x0F,
}

func (kb *Keyboard) Start() {
	kb.keysDown = make(map[byte]time.Time)
	go kb.listner()
}

func (kb *Keyboard) Stop() {
	keyboard.Close()
}

// DoKeyEventUpdates tracks key presses and releases. It should be called 
// on every exection cycle to periodically keep track of keys.
func (kb *Keyboard) DoKeyEventUpdates() {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	kb.prevKeysPressed = kb.currentKeysPressed
	kb.currentKeysPressed = kb.tempKeysPressed
}

func (kb *Keyboard) JustPressed(key byte) bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	return kb.currentKeysPressed[key] && !kb.prevKeysPressed[key]
}

func (kb *Keyboard) JustReleased(key byte) bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	return !kb.currentKeysPressed[key] && kb.prevKeysPressed[key]
}

func (kb *Keyboard) IsPressed(key byte) bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	return kb.currentKeysPressed[key]
}

func (kb *Keyboard) listner() {
	// Create a channel to poll for key inputs
	keysEvents, err := keyboard.GetKeys(keyChannelSize)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case event := <-keysEvents: {
			pressedChar := strings.ToLower(string(event.Rune))
			if event.Key == keyboard.KeyCtrlC {
				fmt.Printf("Press <Ctrl-c> once again to exit!\n")
				kb.Stop()
				return	
			}
			if charIdx, ok := KeyMap[pressedChar]; ok {
				kb.mu.Lock()
				kb.keysDown[charIdx] = time.Now()
				kb.tempKeysPressed[charIdx] = true
				kb.mu.Unlock()
			}
		}
		default: {
			now := time.Now()
			for k, t := range kb.keysDown {
				if now.Sub(t) >= releaseDelay {
					kb.mu.Lock()
					// Mark key as released
					delete(kb.keysDown, k)
					kb.tempKeysPressed[k] = false
					kb.mu.Unlock()
				}
			}
		}
		}
	}
}
