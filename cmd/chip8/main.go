package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/beniceok/chip8"
	"github.com/beniceok/chip8/components"
	"github.com/beniceok/chip8/debug"
	"github.com/davecgh/go-spew/spew"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	file := flag.String("file", "", "enter a relative file path to load")
	dump := flag.Bool("dump", false, "dump the rom contents")
	debugEnabled := flag.Bool("debug", false, "debug output")
	scale := flag.Int("scale", 4, "scale each display pixel")
	delay := flag.Int("delay", 0, "delay between each cpu cycle in milliseconds, for debugging")
	_ = flag.Bool("sound", false, "enable sound, false by default, TODO: not implemented")

	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(0)
	}

	c8 := chip8.NewChip8()
	debugger := debug.Debugger{Chip8: c8}

	c8.LoadRom(*file)

	if *dump {
		spew.Dump(c8.GetMem())
	}

	quit := make(chan bool, 1)

	if *debugEnabled {
		debugger.Start(quit)
	}

	rand.Seed(time.Now().Unix())

	pixelgl.Run(
		components.MainLoop(
			c8,
			float64(*scale),
			float64(chip8.Width),
			float64(chip8.Height),
			*delay,
			quit))
}
