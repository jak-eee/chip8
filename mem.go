package chip8

// Start is the memory address the program counter first points to.
const Start = 0x0200

// nextOpcode gets the next opcode from memory.
func (c *Chip8) nextOpcode() uint16 {
	var oc uint16
	oc = uint16(c.mem[c.pgCtr])
	oc <<= 8
	oc |= uint16(c.mem[c.pgCtr+1])
	c.pgCtr += 2

	return oc
}

// clearMem clears the registers, stack, and memory.
func (c *Chip8) clearMem() {
	c.mem = [4096]byte{}
	c.reg = [16]byte{}

	for {
		_, err := c.stack.Pop()
		if err != nil {
			break
		}
	}
}

// GetMem returns the memory.
func (c *Chip8) GetMem() [4096]byte {
	return c.mem
}
