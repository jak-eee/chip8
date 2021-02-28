package components

import (
	"time"

	"github.com/beniceok/chip8"
)

// StartSpeaker enables sound.
func StartSpeaker() {
	go func() {
		for {
			if chip8.Beep {
				{
					// TODO: implement.
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
