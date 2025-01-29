package chip8

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// LoadRom loads the given file into memory.
func (c *Chip8) LoadRom(path string) {
	c.loadFont()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("could not get current directory: ", err)
	}

	fullPath := dir + "/" + path

	bs, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Fatal("error reading file: ", err)
	}

	c.romName = filepath.Base(path)

	for i, b := range bs {
		c.mem[Start+i] = byte(b)
	}
}
