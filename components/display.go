package components

import (
	"image/color"
	_ "image/png"
	"log"
	"time"

	"github.com/jak-eee/chip8"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// MainLoop is the main window loop provided by the pixelgl graphics opengl wrapper.
// scale is a multiplier for drawing raw pixels and defines the window size, it must be a positive integer.
func MainLoop(
	chip8 *chip8.Chip8,
	scale float64,
	width float64,
	height float64,
	delay int,
	quit chan bool,
) func() {
	return func() {
		title := chip8.GetRomName()

		cfg := pixelgl.WindowConfig{
			Title:  title,
			Bounds: pixel.R(0, 0, width*scale, height*scale),
			VSync:  true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}

		prevScale := scale

		for !win.Closed() {
			ts := time.Now()
			if chip8.IsAwaitingInput() {
				title = "(press any key)"
				win.SetTitle(title)
			} else if title != chip8.GetRomName() {
				title = chip8.GetRomName()
				win.SetTitle(title)
			}

			scale = handleKeyboard(chip8, win, scale)

			// Realtime scaling simply destroys and recreates the window, this is buggy.
			if scale != prevScale {
				prevScale = scale
				cfg := pixelgl.WindowConfig{
					Title:  title,
					Bounds: pixel.R(0, 0, width*scale, height*scale),
					VSync:  true,
				}
				win.Destroy()
				win, err = pixelgl.NewWindow(cfg)
				if err != nil {
					log.Fatal("error resizing", err)
				}
			}

			if !chip8.Halted() {
				chip8.FetchDecodeCycle()
			}

			if !chip8.Halted() && chip8.Drawn() {
				imd := imdraw.New(nil)
				buffer := chip8.GetBuffer()
				imd.Reset()
				imd.Clear()
				imd.Color = color.RGBA{200, 220, 255, 255}

				for y := 0; y < len(buffer); y++ {
					for x := 0; x < len(buffer[y]); x++ {
						if buffer[y][x] == 0x01 {
							imd.Push(pixel.V((float64(x))*scale, height*scale-(float64(y+1))*scale))
							imd.Push(pixel.V((float64(x))*scale+scale, height*scale-(float64(y+1))*scale+scale))
							imd.Rectangle(0)
						}
					}
				}

				win.Clear(color.RGBA{20, 30, 40, 100})
				imd.Draw(win)
			}

			win.Update()

			if !chip8.Halted() && chip8.Drawn() {
				chip8.RecordFrame(time.Since(ts))
			}

			if delay > 0 {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
		quit <- true
	}
}
