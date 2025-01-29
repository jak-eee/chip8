package chip8

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(0)
}

func Test_OpCodeANNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xAF12)
	c8.cycle(oc)
	if c8.addrI != 0x0F12 {
		t.Logf("c8.addrI was %x, expected F12", c8.addrI)
		t.Fail()
	}
}

func Test_OpCodeBNNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xBF12)
	c8.cycle(oc)
	if c8.pgCtr != 0x0F12 {
		t.Logf("c8.pgCtr was %x, expected 0x0F12", c8.addrI)
		t.Fail()
	}
}

func Test_OpCode1NNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x1A05)
	c8.cycle(oc)
	if c8.pgCtr != 0xA05 {
		t.Logf("c8.pgCtr was %x, expected A05", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCode2NNN(t *testing.T) {
	c8 := NewChip8()
	c8.pgCtr = 0x0123
	oc := uint16(0x2A05)
	c8.cycle(oc)
	v, _ := c8.stack.Pop()
	if v != 0x0123 {
		t.Logf("last word in stack was %x, expected 0x0123", c8.pgCtr)
		t.Fail()
	}
	if c8.pgCtr != 0x0A05 {
		t.Logf("c8.pgCtr was %x, expected 0x0A05", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCode3XNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x3F8C)
	c8.reg[0xF] = 0x8C
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 202", c8.pgCtr)
		t.Fail()
	}
	c8.reg[0xF] = 0x11
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 202", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCode4XNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x4F8C)
	c8.reg[0xF] = 0x11
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 2", c8.pgCtr)
		t.Fail()
	}
	c8.reg[0xF] = 0x8C
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 2", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCode5XY0(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x5AB0)
	c8.reg[0xA] = 0b10
	c8.reg[0xB] = 0b11
	c8.cycle(oc)
	if c8.pgCtr != 0x200 {
		t.Logf("c8.pgCtr was %x, expected 200", c8.pgCtr)
		t.Fail()
	}
	c8.reg[0xA] = 0b11
	c8.reg[0xB] = 0b11
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 202", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCode00EE(t *testing.T) {
	c8 := NewChip8()
	c8.pgCtr = 0x0123
	oc := uint16(0x2A05)
	c8.cycle(oc)
	oc = uint16(0x00EE)
	c8.cycle(oc)
	if c8.pgCtr != 0x0123 {
		t.Logf("c8.pgCtr was %x, expected 0x0123", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCode6XNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x6AB2)
	c8.cycle(oc)
	if c8.reg[0xA] != 0xB2 {
		t.Logf("c8.reg[0xA] was %x, expected b2", c8.reg[0xA])
		t.Fail()
	}
}

func Test_OpCode7XNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x7514)
	c8.reg[0x5] = 0x02
	c8.cycle(oc)
	if c8.reg[0x5] != 0x16 {
		t.Logf("c8.reg[0x5] was %x, expected 0x16", c8.reg[0x5])
		t.Fail()
	}
}

func Test_OpCode8XY0(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB0)
	c8.reg[0xB] = 0x02
	c8.cycle(oc)
	if c8.reg[0xA] != 0x02 {
		t.Logf("c8.reg[0xA] was %x, expected 0x02", c8.reg[0xA])
		t.Fail()
	}
}

func Test_OpCode8XY1(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB1)
	c8.reg[0xA] = 0xA1
	c8.reg[0xB] = 0xFA
	c8.cycle(oc)
	if c8.reg[0xA] != 0xFB {
		t.Logf("c8.reg[0xA] was %x, expected 0xFB", c8.reg[0xA])
		t.Fail()
	}
}

func Test_OpCode8XY2(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB2)
	c8.reg[0xA] = 0xA1
	c8.reg[0xB] = 0xFA
	c8.cycle(oc)
	if c8.reg[0xA] != 0xA0 {
		t.Logf("c8.reg[0xA] was %x, expected 0xA0", c8.reg[0xA])
		t.Fail()
	}
}

func Test_OpCode8XY3(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB3)
	c8.reg[0xA] = 0xA1
	c8.reg[0xB] = 0xFA
	c8.cycle(oc)
	if c8.reg[0xA] != 0x5B {
		t.Logf("c8.reg[0xA] was %x, expected 0x5B", c8.reg[0xA])
		t.Fail()
	}
}

func Test_OpCode8XY4(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB4)
	c8.reg[0xA] = 0xA1
	c8.reg[0xB] = 0xFA
	c8.cycle(oc)
	if c8.reg[0xA] != 0x9B {
		t.Logf("c8.reg[0xA] was %x, expected 0x9B", c8.reg[0xA])
		t.Fail()
	}
	if c8.reg[0xF] != 0x1 {
		t.Logf("c8.reg[0xF] was %x, expected 1", c8.reg[0xF])
		t.Fail()
	}

	c8.reg[0xA] = 0x1
	c8.reg[0xB] = 0x1
	c8.cycle(oc)
	if c8.reg[0xA] != 0x2 {
		t.Logf("c8.reg[0xA] was %x, expected 0x02", c8.reg[0xA])
		t.Fail()
	}
	if c8.reg[0xF] != 0x0 {
		t.Logf("c8.reg[0xF] was %x, expected 0", c8.reg[0xF])
		t.Fail()
	}
}

func Test_OpCode8XY5(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB5)
	c8.reg[0xA] = 0x11
	c8.reg[0xB] = 0xFF
	c8.cycle(oc)
	if c8.reg[0xA] != 0x12 {
		t.Logf("c8.reg[0xA] was %x, expected 0x12", c8.reg[0xA])
		t.Fail()
	}
	if c8.reg[0xF] != 0x0 {
		t.Logf("c8.reg[0xF] was %x, expected 0", c8.reg[0xF])
		t.Fail()
	}

	c8.reg[0xA] = 0x1
	c8.reg[0xB] = 0x1
	c8.cycle(oc)
	if c8.reg[0xA] != 0x0 {
		t.Logf("c8.reg[0xA] was %x, expected 0x0", c8.reg[0xA])
		t.Fail()
	}
	if c8.reg[0xF] != 0x1 {
		t.Logf("c8.reg[0xF] was %x, expected 1", c8.reg[0xF])
		t.Fail()
	}
}

func Test_OpCode8XY6(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB6)
	c8.reg[0xA] = 0b10010111
	c8.cycle(oc)
	if c8.reg[0xF] != 1 {
		t.Logf("c8.reg[0xF] was %x, expected 1", c8.reg[0xF])
		t.Fail()
	}

	c8.reg[0xA] = 0b10010110
	c8.cycle(oc)
	if c8.reg[0xF] != 0 {
		t.Logf("c8.reg[0xF] was %x, expected 0", c8.reg[0xF])
		t.Fail()
	}
}

func Test_OpCode8XY7(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8AB7)
	c8.reg[0xA] = 0x11
	c8.reg[0xB] = 0xFF
	c8.cycle(oc)
	if c8.reg[0xA] != 0xEE {
		t.Logf("c8.reg[0xA] was %x, expected 0xEE", c8.reg[0xA])
		t.Fail()
	}
	if c8.reg[0xF] != 0x1 {
		t.Logf("c8.reg[0xF] was %x, expected 1", c8.reg[0xF])
		t.Fail()
	}

	c8.reg[0xA] = 0xFF
	c8.reg[0xB] = 0x11
	c8.cycle(oc)
	if c8.reg[0xA] != 0x12 {
		t.Logf("c8.reg[0xA] was %x, expected 0x12", c8.reg[0xA])
		t.Fail()
	}
	if c8.reg[0xF] != 0x0 {
		t.Logf("c8.reg[0xF] was %x, expected 0", c8.reg[0xF])
		t.Fail()
	}
}

func Test_OpCode8XYE(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x8ABE)
	c8.reg[0xA] = 0b10010111
	c8.cycle(oc)
	if c8.reg[0xF] != 1 {
		t.Logf("c8.reg[0xF] was %x, expected 1", c8.reg[0xF])
		t.Fail()
	}

	c8.reg[0xA] = 0b00010111
	c8.cycle(oc)
	if c8.reg[0xF] != 0 {
		t.Logf("c8.reg[0xF] was %x, expected 0", c8.reg[0xF])
		t.Fail()
	}
}

func Test_OpCode9XY0(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0x9AB0)
	c8.reg[0xA] = 0x1
	c8.reg[0xB] = 0x1
	c8.cycle(oc)
	if c8.pgCtr != 0x200 {
		t.Logf("c8.pgCtr was %x, expected 200", c8.pgCtr)
		t.Fail()
	}
	c8.reg[0xA] = 0x0
	c8.reg[0xB] = 0x1
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 202", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCodeCXNN(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xC4A1)
	c8.reg[4] = 0x00
	c8.cycle(oc)
	if c8.reg[0x4] != 0x81 {
		t.Logf("c8.reg[0x4] was %x, expected 81", c8.reg[0x4])
		t.Fail()
	}
}

func Test_OpCodeEX9E(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xEF9E)
	c8.reg[0xF] = KeyF
	c8.keyPressed[KeyF] = true
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 202", c8.pgCtr)
		t.Fail()
	}
}

func Test_OpCodeEXA1(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xEFA1)
	c8.reg[0xF] = Key0
	c8.keyPressed[Key0] = false
	c8.cycle(oc)
	if c8.pgCtr != 0x202 {
		t.Logf("c8.pgCtr was %x, expected 202", c8.pgCtr)
		t.Fail()
	}
}

func TestOpCodeFX07(t *testing.T) {
	c8 := NewChip8()
	c8.reset()
	oc := uint16(0xF107)
	c8.delayTimer = 0x12
	c8.cycle(oc)
	if c8.reg[0x1] != 0x12 {
		t.Logf("c8.reg[0x1] was %x, expected 0x12", c8.reg[0x1])
		t.Fail()
	}
}

func TestOpCodeFX15(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xF115)
	c8.reg[0x1] = 0xAA
	c8.cycle(oc)
	if c8.delayTimer != 0xAA {
		t.Logf("DelayTimer was %x, expected 0xAA", c8.delayTimer)
		t.Fail()
	}
}

func TestOpCodeFX18(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xF118)
	c8.reg[0x1] = 0xAA
	c8.cycle(oc)
	if c8.soundTimer != 0xAA {
		t.Logf("SoundTimer was %x, expected 0xAA", c8.soundTimer)
		t.Fail()
	}
}

func TestOpCodeFX0A(t *testing.T) {
	c8 := NewChip8()
	done := make(chan bool, 1)
	timer := time.NewTimer(time.Millisecond * 10)
	oc := uint16(0xF90A)
	c8.cycle(oc)
	go func() {
		c8.awaitInput <- Key1
		<-timer.C
		done <- true
	}()
	<-done
	if c8.reg[0x9] != Key1 {
		t.Logf("c8.reg[0x9] was %v, expected 0x1", c8.reg[0x9])
		t.Fail()
	}
}

func TestOpCodeFX33(t *testing.T) {
	c8 := NewChip8()
	oc := uint16(0xF633)
	c8.reg[0x6] = 241
	c8.cycle(oc)

	if c8.mem[0] != 2 {
		t.Logf("Mem[0] was %x, expected 2", c8.mem[0])
		t.Fail()
	}
	if c8.mem[1] != 4 {
		t.Logf("Mem[1] was %x, expected 4", c8.mem[0])
		t.Fail()
	}
	if c8.mem[2] != 1 {
		t.Logf("Mem[2] was %x, expected 1", c8.mem[0])
		t.Fail()
	}
}
