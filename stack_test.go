package chip8

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestStack(t *testing.T) {
	s := new(stack)

	for i := 0; i <= 10; i++ {
		s.Push(uint16(i))
	}

	for i := 0; i <= 10; i++ {
		v, err := s.Pop()
		exp := 10 - i
		if v != uint16(exp) {
			t.Logf("expected %x, got value %x", exp, v)
			t.Fail()
		}

		if v == uint16(0) {
			spew.Dump(err)
		}
	}

	_, err := s.Pop()
	if err.Error() != "stack is empty" {
		t.Logf("expected error \"stack is empty\", got error %s", err)
		t.Fail()
	}

	var nilStack *stack
	_, errNil := nilStack.Pop()
	if errNil.Error() != "stack is nil" {
		t.Logf("expected error \"stack is nil\", got error %s", errNil)
		t.Fail()
	}
}
