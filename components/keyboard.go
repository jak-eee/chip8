package components

import (
	"github.com/jak-eee/chip8"
	"github.com/faiface/pixel/pixelgl"
)

// keyMap is the mapping of emulated key values to real keys.
// Edit this map to change the keyboard bindings.
var keyMap = map[uint8]pixelgl.Button{
	chip8.Key0: pixelgl.Key1,
	chip8.Key1: pixelgl.Key2,
	chip8.Key2: pixelgl.Key3,
	chip8.Key3: pixelgl.Key4,
	chip8.Key4: pixelgl.KeyQ,
	chip8.Key5: pixelgl.KeyW,
	chip8.Key6: pixelgl.KeyE,
	chip8.Key7: pixelgl.KeyR,
	chip8.Key8: pixelgl.KeyA,
	chip8.Key9: pixelgl.KeyS,
	chip8.KeyA: pixelgl.KeyD,
	chip8.KeyB: pixelgl.KeyF,
	chip8.KeyC: pixelgl.KeyZ,
	chip8.KeyD: pixelgl.KeyX,
	chip8.KeyE: pixelgl.KeyC,
	chip8.KeyF: pixelgl.KeyV,
}

// handleKeyboard sets the KeyPressed slice if a real key is pressed.
func handleKeyboard(c8 *chip8.Chip8, win *pixelgl.Window, scale float64) float64 {
	for key, realKey := range keyMap {
		if win.Pressed(realKey) {
			c8.KeyPressed(key, true)
		} else {
			c8.KeyPressed(key, false)
		}
		if win.JustReleased(realKey) {
			c8.KeyPressed(key, false)
		}
	}

	return handleScale(win, scale)
}

// handleScale handles the scaling keys not part of the interpreter.
func handleScale(win *pixelgl.Window, scale float64) float64 {
	if win.JustPressed(pixelgl.KeyMinus) {
		if scale > 2 {
			scale--
		}
	}
	if win.JustPressed(pixelgl.KeyEqual) {
		scale++
	}

	return scale
}
