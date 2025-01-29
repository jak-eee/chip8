package chip8

// Hardcoded screen width.
const Width = 64

// Hardcoded screen height.
const Height = 32

// displayClear 0's out the display buffer.
func (c *Chip8) displayClear() {
	for y := 0; y < len(c.buffer); y++ {
		for x := 0; x < len(c.buffer[0]); x++ {
			c.buffer[y][x] = 0x0
		}
	}
	c.drawn = true
}

// draw draws a sprite to the display buffer from memory and returns whether or not there
// is a collision with an existing pixel in the buffer.
func (c *Chip8) draw(x byte, y byte, width byte, height byte) bool {
	c.drawn = true
	flip := false
	var curY int = int(y)
	var curX int = int(x) + int(width) - 1

	for i := c.addrI; i < c.addrI+uint16(height); i++ {
		b := c.mem[i]
		for j := 0; j < 8; j++ {
			p := b >> (j) & 0b1
			if curY > 31 {
				curY = curY % 32
			}
			if curX > 63 {
				curX = curX % 64
			}

			if p == 0x1 && c.buffer[curY][curX] == uint8(p) {
				flip = true
			}

			c.buffer[curY][curX] ^= uint8(p)

			curX -= 1
			if curX < 0 {
				curX = 63
			}

			if j == 7 {
				curX = int(x) + int(width) - 1
				curY += 1
			}
		}
	}

	return flip
}
