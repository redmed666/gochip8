package chip8

import (
	"os"

	"github.com/redmed666/gochip8/util"
)

// Chip8 structure
type Chip8 struct {
	opcode     uint16      // 1 opcode == 2 bytes
	memory     [4096]uint8 // memory == 4096 bytes
	V          [16]uint8   // array of registers => 1 register == 1 byte
	I          uint16
	PC         uint16
	gfx        [64 * 32]uint8 // black and white screen of 64 pixels by 32
	delayTimer uint8
	soundTimer uint8
	stack      [16]uint16 // chip8 doesn't have a stack but it will be useful
	SP         uint16
	key        [16]uint8 // keypad
	drawFlag   byte
}

var fontset = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func (chip8 *Chip8) setupGraphics() {
}

func (chip8 *Chip8) setupInput() {

}

func (chip8 *Chip8) initialize() {
	chip8.PC = 0x200
	chip8.opcode = 0
	chip8.I = 0
	chip8.SP = 0

	for i := 0; i < 64*32; i++ {
		chip8.gfx[i] = 0
	}

	for i := 0; i < 16; i++ {
		chip8.stack[i] = 0
		chip8.V[i] = 0
	}

	for i := 0; i < 4096; i++ {
		chip8.memory[i] = 0
	}

	for i := 0; i < 80; i++ {
		chip8.memory[i] = fontset[i]
	}

	chip8.delayTimer = 0
	chip8.soundTimer = 0
}

func (chip8 *Chip8) loadGame(gamePath string) {
	file, err := os.Open(gamePath)
	util.CheckError(err)
	fileStat, err := file.Stat()
	util.CheckError(err)
	data := make([]byte, fileStat.Size())
	file.Read(data)

	for i := 0; i < len(data); i++ {
		chip8.memory[512+i] = data[i] // Start filling memory at location 0x200 == 512 in the chip8
	}
}

func (chip8 *Chip8) emulateCycle() {

}

func (chip8 *Chip8) setKeys() {

}
