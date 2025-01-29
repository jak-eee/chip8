package chip8

// The 16 available key constants
const Key0 = 0x0
const Key1 = 0x1
const Key2 = 0x2
const Key3 = 0x3
const Key4 = 0x4
const Key5 = 0x5
const Key6 = 0x6
const Key7 = 0x7
const Key8 = 0x8
const Key9 = 0x9
const KeyA = 0xA
const KeyB = 0xB
const KeyC = 0xC
const KeyD = 0xD
const KeyE = 0xE
const KeyF = 0xF

// KeyPressed sets the key pressed status and sends to the
// awaitingInput channel if program is halted awaiting input.
func (c *Chip8) KeyPressed(key byte, pressed bool) {
	c.keyPressed[key] = pressed
	if c.awaitingInput && pressed {
		c.awaitInput <- key
	}
}

// IsAwaitingInput is true if the program is expecting input.
func (c *Chip8) IsAwaitingInput() bool {
	return c.awaitingInput
}
