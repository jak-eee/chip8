package chip8

import "errors"

// stack is non-thread safe uint16 stack.
// It doesn't need to be thread safe as there is only a single emulated "thread".
type stack []uint16

// Push, not nil safe
func (s *stack) Push(v uint16) {
	*s = append(*s, v)
}

// Pop, empty or nil returns error
func (s *stack) Pop() (uint16, error) {
	if s == nil {
		return 0, errors.New("stack is nil")
	}
	if len(*s) == 0 {
		return 0, errors.New("stack is empty")
	}

	v := (*s)[len(*s)-1]

	if len(*s) > 0 {
		*s = (*s)[0 : len(*s)-1]
	}

	return v, nil
}
