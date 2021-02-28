// chip8 contains all functions related to the chip-8 interpreter; cpu input, timers, sound, and graphics.
// See the components subpackage for the external implementations of the window, display, keyboard, and sound.
// The components package directly references the chip8 package. The chip8 package doesn't reference the components package.
// Implementing a new window manager, display routine, or changing keyboard bindings involves editing the components packages.
// Changing the interpreter involves editing the chip8 package.
// The comments above the opcode functions (OpCodeXXX) are taken directly from the wikipedia article https://en.wikipedia.org/wiki/CHIP-8#cite_note-increment-17.
package chip8

import (
	"log"
	"math/rand"
	"time"
)

// Expected cycle time in nanoseconds.
var CycleTimeNS float64

// Cycles.
// TODO: make part of type and make configurable.
const Cycles = 60

// OpCodesPerCycle is the number of op codes to decode and execute per frame.
// TODO: make part of type and make configurable.
const OpCodesPerCycle = 7

type Chip8 struct {
	// mem is 4 kilobytes.
	mem [4096]byte

	// reg is 16 8-bit registers.
	reg [16]byte

	// addrI is the current memory address.
	addrI uint16

	// pgCtr is the program counter.
	pgCtr uint16

	// stack is a stack of uint16s.
	stack stack

	// buffer is the current display buffer.
	buffer [Height][Width]uint8

	// drawn signifies that something has been written to the buffer.
	drawn bool

	// delayTimer counts down to zero 60 times per "second".
	// (Really once per cycle at roughly 60 cycles per second).
	delayTimer byte

	// soundTimer counts down to zero 60 times per "second".
	// (Really once per cycle at roughly 60 cycles per second).
	soundTimer byte

	// keyPressed is a slice of bools that correspond to whether or not an emulated key is pressed
	keyPressed [16]bool

	// awaitInput is set if the system is halted awaiting a key press.
	awaitInput chan byte

	// awaitingInput is set to true if the program is expecting a key to be sent to AwaitInput.
	awaitingInput bool

	// halt is whether or not interpreter is halted.
	halt bool

	romName string

	// info for debugging only.
	framesI    uint64
	framesAvg  float64
	frameTimes []time.Duration
	opCodesI   uint64
	opCodesAvg float64
	opCodes    []uint16

	tz *time.Time
}

func NewChip8() *Chip8 {
	c8 := &Chip8{pgCtr: Start}
	c8.awaitInput = make(chan byte, 1)
	c8.reset()

	return c8
}

func (c *Chip8) RecordFrame(ft time.Duration) {
	c.framesI++
	c.framesAvg = float64(c.framesI) / float64(time.Since(*c.tz).Seconds())
	c.frameTimes = append(c.frameTimes, ft)
	if len(c.frameTimes) > 64 {
		c.frameTimes = c.frameTimes[len(c.frameTimes)-64:]
	}
}

// Functions provided for debugging.
func (c *Chip8) GetFramesI() uint64 {
	return c.framesI
}

func (c *Chip8) GetFrameTimes() []time.Duration {
	return c.frameTimes
}

func (c *Chip8) GetFramesAvg() float64 {
	return c.framesAvg
}

func (c *Chip8) GetOpCodes() []uint16 {
	return c.opCodes
}

func (c *Chip8) GetOpCodesI() uint64 {
	return c.opCodesI
}

func (c *Chip8) GetOpCodesAvg() float64 {
	return float64(c.opCodesAvg)
}

func (c *Chip8) GetDelayTimer() byte {
	return c.delayTimer
}

func (c *Chip8) GetSoundTimer() byte {
	return c.soundTimer
}

func (c *Chip8) GetReg() [16]byte {
	return c.reg
}

func (c *Chip8) GetKeyPressed() [16]bool {
	return c.keyPressed
}

func init() {
	// set the expected time per frame.
	CycleTimeNS = 1.0 / float64(Cycles) * 1000000000
}

func (c *Chip8) GetRomName() string {
	return c.romName
}

func (c *Chip8) Halted() bool {
	return c.halt
}

func (c *Chip8) Drawn() bool {
	return c.drawn
}

// reset resets the program counter, resets the current address, clears the stack, and clears the registers.
func (c *Chip8) reset() {
	c.delayTimer = 0x0
	c.soundTimer = 0x0
	c.pgCtr = Start
	c.addrI = 0
	tz := time.Now()
	c.tz = &tz
	c.clearMem()
}

func (c *Chip8) GetBuffer() [Height][Width]uint8 {
	return c.buffer
}

// FetchDecodeCycle is the main program loop.
func (c *Chip8) FetchDecodeCycle() {
	c.drawn = false
	if !c.halt {
		c.decrementTimers()
		for i := 0; i < OpCodesPerCycle; i++ {
			c.cycle(c.nextOpcode())
		}
	}
}

func (c *Chip8) decrementTimers() {
	if c.soundTimer > 0 {
		Beep = true
		c.soundTimer--
	} else {
		Beep = false
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}
}

// cycle decodes the opcode and executes the corresponding function.
func (c *Chip8) cycle(oc uint16) {
	c.opCodesI++
	c.opCodes = append(c.opCodes, oc)
	if len(c.opCodes) > 16 {
		c.opCodes = c.opCodes[len(c.opCodes)-16:]
	}
	c.opCodesAvg = float64(c.opCodesI) / float64(time.Since(*c.tz).Seconds())

	switch oc & 0xF000 {
	case 0x0000:
		c.decode0x0000(oc)
	case 0x1000:
		c.opCode1NNN(oc)
	case 0x2000:
		c.opCode2NNN(oc)
	case 0x3000:
		c.opCode3XNN(oc)
	case 0x4000:
		c.opCode4XNN(oc)
	case 0x5000:
		c.opCode5XY0(oc)
	case 0xC000:
		c.opCodeCXNN(oc)
	case 0xD000:
		c.opCodeDXYN(oc)
	case 0x6000:
		c.opCode6XNN(oc)
	case 0x7000:
		c.opCode7XNN(oc)
	case 0x8000:
		c.decode0x8000(oc)
	case 0x9000:
		c.opCode9XY0(oc)
	case 0xA000:
		c.opCodeANNN(oc)
	case 0xB000:
		c.opCodeBNNN(oc)
	case 0xE000:
		c.decode0xE000(oc)
	case 0xF000:
		c.decode0xF000(oc)
	default:
		log.Printf("unknown opcode 0x%x\n", oc)
	}
}

func (c *Chip8) decode0x8000(oc uint16) {
	switch oc & 0x000F {
	case 0x0:
		c.opCode8XY0(oc)
	case 0x1:
		c.opCode8XY1(oc)
	case 0x2:
		c.opCode8XY2(oc)
	case 0x3:
		c.opCode8XY3(oc)
	case 0x4:
		c.opCode8XY4(oc)
	case 0x5:
		c.opCode8XY5(oc)
	case 0x6:
		c.opCode8XY6(oc)
	case 0x7:
		c.opCode8XY7(oc)
	case 0xE:
		c.opCode8XYE(oc)
	}
}

func (c *Chip8) decode0x0000(oc uint16) {
	switch oc & 0x00FF {
	case 0x00E0:
		c.opCode00E0(oc)
	case 0x00EE:
		c.opCode00EE(oc)
	default:
		c.opCode0NNN(oc)
	}
}

func (c *Chip8) decode0xE000(oc uint16) {
	switch oc & 0x00FF {
	case 0x9E:
		c.opCodeEX9E(oc)
	case 0xA1:
		c.opCodeEXA1(oc)
	}
}

func (c *Chip8) decode0xF000(oc uint16) {
	switch oc & 0x00FF {
	case 0x07:
		c.opCodeFX07(oc)
	case 0x15:
		c.opCodeFX15(oc)
	case 0x18:
		c.opCodeFX18(oc)
	case 0x1E:
		c.opCodeFX1E(oc)
	case 0x0A:
		c.opCodeFX0A(oc)
	case 0x29:
		c.opCodeFX29(oc)
	case 0x33:
		c.opCodeFX33(oc)
	case 0x55:
		c.opCodeFX55(oc)
	case 0x65:
		c.opCodeFX65(oc)
	}
}

// Call
// Calls machine code routine (RCA 1802 for COSMAC VIP) at address NNN.
// Not necessary for most ROMs.
func (c *Chip8) opCode0NNN(oc uint16) {
	// no op
}

// Jumps to address NNN.
func (c *Chip8) opCode1NNN(oc uint16) {
	nnn := oc & 0x0FFF

	c.pgCtr = nnn
}

// Calls subroutine at NNN.
func (c *Chip8) opCode2NNN(oc uint16) {
	nnn := oc & 0x0FFF

	c.stack.Push(c.pgCtr)

	c.pgCtr = nnn
}

func (c *Chip8) opCode00E0(oc uint16) {
	c.displayClear()
}

// Returns from a subroutine.
func (c *Chip8) opCode00EE(oc uint16) {
	v, err := c.stack.Pop()

	if err == nil {
		c.pgCtr = v
	} else {
		log.Fatal(err)
	}
}

// if(Vx==NN) Skips the next instruction if VX equals NN.
func (c *Chip8) opCode3XNN(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	nn := oc & 0x00FF

	if byte(nn) == c.reg[x] {
		c.pgCtr += 2
	}
}

// if(Vx!=NN) Skips the next instruction if VX doesn't equal NN.
func (c *Chip8) opCode4XNN(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	nn := oc & 0x00FF

	if byte(nn) != c.reg[x] {
		c.pgCtr += 2
	}
}

// if(Vx==Vy) Skips the next instruction if VX equals VY.
func (c *Chip8) opCode5XY0(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	if c.reg[x] == c.reg[y] {
		c.pgCtr += 2
	}
}

// 6XNN Const Vx = NN Sets VX to NN.
func (c *Chip8) opCode6XNN(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	nn := oc & 0x00FF

	c.reg[x] = byte(nn)
}

// 7XNN Const Vx += NN Adds NN to VX. (Carry flag is not changed)
func (c *Chip8) opCode7XNN(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	nn := oc & 0x00FF

	c.reg[x] = c.reg[x] + byte(nn)
}

// 8XY0 Assign Vx=Vy Sets VX to the value of VY.
func (c *Chip8) opCode8XY0(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	c.reg[x] = c.reg[y]
}

// 8XY1 BitOp Vx=Vx|Vy Sets VX to VX or VY. (Bitwise OR operation)
func (c *Chip8) opCode8XY1(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	r := c.reg[x] | c.reg[y]

	c.reg[x] = r
}

// 8XY2 BitOp Vx=Vx&Vy Sets VX to VX and VY. (Bitwise AND operation)
func (c *Chip8) opCode8XY2(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	r := c.reg[x] & c.reg[y]

	c.reg[x] = r
}

// 8XY3[a] BitOp Vx=Vx^Vy Sets VX to VX xor VY.
func (c *Chip8) opCode8XY3(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	r := c.reg[x] ^ c.reg[y]

	c.reg[x] = r
}

// 8XY4 Math
// Vx += Vy Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
func (c *Chip8) opCode8XY4(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	var res int = int(c.reg[x]) + int(c.reg[y])

	if res > 255 {
		c.reg[0xF] = 1
	} else {
		c.reg[0xF] = 0
	}

	c.reg[x] = byte(res)
}

// 8XY5 Math
// Vx -= Vy VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
func (c *Chip8) opCode8XY5(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	if c.reg[y] > c.reg[x] {
		c.reg[0xF] = 0
	} else {
		c.reg[0xF] = 1
	}

	c.reg[x] = c.reg[x] - c.reg[y]
}

// 8XY6[a]
// BitOp Vx>>=1 Stores the least significant bit of VX in VF and then shifts VX to the right by 1.[b]
func (c *Chip8) opCode8XY6(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	lsvx := c.reg[x] & 0b00000001

	c.reg[0xF] = lsvx

	c.reg[x] = c.reg[x] >> 1
}

// 8XY7[a] Math
// Vx=Vy-Vx Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
func (c *Chip8) opCode8XY7(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	if c.reg[y] < c.reg[x] {
		c.reg[0xF] = 0
	} else {
		c.reg[0xF] = 1
	}

	c.reg[x] = c.reg[y] - c.reg[x]
}

// 8XYE[a]
// BitOp Vx<<=1 Stores the most significant bit of VX in VF and then shifts VX to the left by 1.[b]
func (c *Chip8) opCode8XYE(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	lsvx := c.reg[x] & 0b10000000

	if lsvx>>7 == 1 {
		c.reg[0xF] = 1
	} else {
		c.reg[0xF] = 0
	}

	c.reg[x] = c.reg[x] << 1
}

// 9XY0 Cond
// if(Vx!=Vy)
// Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
func (c *Chip8) opCode9XY0(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	if c.reg[x] != c.reg[y] {
		c.pgCtr += 2
	}
}

// ANNN MEM I = NNN Sets I to the address NNN.
func (c *Chip8) opCodeANNN(oc uint16) {
	nnn := oc & 0x0FFF

	c.addrI = nnn
}

// BNNN Flow PC=V0+NNN Jumps to the address NNN plus V0.
func (c *Chip8) opCodeBNNN(oc uint16) {
	nnn := oc & 0x0FFF

	c.pgCtr = nnn + uint16(c.reg[0x0])
}

// CXNN Rand Vx=rand()&NN
// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
func (c *Chip8) opCodeCXNN(oc uint16) {
	nn := oc & 0x00FF
	x := oc & 0x0F00
	x = x >> 8

	rnd := byte(rand.Intn(255))

	c.reg[x] = rnd & byte(nn)
}

// Disp draw(Vx,Vy,N)
// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height
// of N+1 pixels. Each row of 8 pixels is read as bit-coded starting from memory
// location I; I value doesn’t change after the execution of this instruction.
// As described above, VF is set to 1 if any screen pixels are flipped from set to
// unset when the sprite is drawn, and to 0 if that doesn’t happen
func (c *Chip8) opCodeDXYN(oc uint16) {
	n := oc & 0x000F
	x := oc & 0x0F00
	x = x >> 8

	y := oc & 0x00F0
	y = y >> 4

	height := n
	flipped := c.draw(c.reg[x], c.reg[y], 8, byte(height))

	if flipped {
		c.reg[0xF] = 0x0001
	} else {
		c.reg[0xF] = 0x0000
	}
}

// EX9E KeyOp if(key()==Vx)
// Skips the next instruction if the key stored in VX is pressed.
// (Usually the next instruction is a jump to skip a code block)
func (c *Chip8) opCodeEX9E(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	if c.keyPressed[c.reg[x]] {
		c.pgCtr += 2
	}
}

// EXA1 KeyOp if(key()!=Vx)
// Skips the next instruction if the key stored in VX isn't pressed.
// (Usually the next instruction is a jump to skip a code block)
func (c *Chip8) opCodeEXA1(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	if !c.keyPressed[c.reg[x]] {
		c.pgCtr += 2
	}
}

// FX07 Timer Vx = get_delay() Sets VX to the value of the delay timer.
func (c *Chip8) opCodeFX07(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	c.reg[x] = c.delayTimer
}

// FX0A	KeyOp Vx = get_key() A key press is awaited, and then stored in VX.
// (Blocking Operation. All instruction halted until next key event)
func (c *Chip8) opCodeFX0A(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	c.halt = true
	c.awaitingInput = true

	go func() {
		key := <-c.awaitInput
		c.reg[x] = key
		c.halt = false
		c.awaitingInput = false
	}()
}

// FX15	Timer delay_timer(Vx) Sets the delay timer to VX.
func (c *Chip8) opCodeFX15(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	c.delayTimer = c.reg[x]
}

// FX18	Sound sound_timer(Vx) Sets the sound timer to VX.
func (c *Chip8) opCodeFX18(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	c.soundTimer = c.reg[x]
}

// FX1E MEM I +=Vx Adds VX to I. VF is not affected.[c]
func (c *Chip8) opCodeFX1E(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	c.addrI += uint16(c.reg[x])
}

// FX29 points c.addrI to font char in reg[x] (0-F)
func (c *Chip8) opCodeFX29(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	c.addrI = uint16(c.reg[x] * 5)
}

/**
FX33 BCD set_BCD(Vx);
*(I+0)=BCD(3);
*(I+1)=BCD(2);
*(I+2)=BCD(1);
Stores the binary-coded decimal representation of VX, with the most significant of three
digits at the address in I, the middle digit at I plus 1, and the least significant
digit at I plus 2. (In other uint16s, take the decimal representation of VX, place the
hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
*/
func (c *Chip8) opCodeFX33(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8
	vx := c.reg[x]

	hundreds := vx / 100
	tens := (vx / 10) % 10
	ones := vx % 10

	c.mem[c.addrI] = hundreds
	c.mem[c.addrI+1] = tens
	c.mem[c.addrI+2] = ones
}

// FX55 MEM reg_dump(Vx,&I)
// Stores V0 to VX (including VX) in memory starting at address I.
// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
func (c *Chip8) opCodeFX55(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	for i := uint16(0); i <= x; i++ {
		c.mem[c.addrI+i] = c.reg[i]
	}

	c.addrI = c.addrI + x + 1
}

// FX65 MEM reg_load(Vx,&I)
// Fills V0 to VX (including VX) with values from memory starting at address I.
// The offset from I is increased by 1 for each value written, but I itself is left unmodified.
func (c *Chip8) opCodeFX65(oc uint16) {
	x := oc & 0x0F00
	x = x >> 8

	for i := uint16(0); i <= x; i++ {
		c.reg[i] = c.mem[c.addrI+i]
	}

	c.addrI = c.addrI + x + 1
}
