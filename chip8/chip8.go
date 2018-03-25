package chip8

import (
	"fmt"
	"os"

	"github.com/redmed666/gochip8/util"
)

// Chip8 structure
type Chip8 struct {
	opcode     uint16         // 1 opcode == 2 bytes
	memory     [4096]uint8    // memory == 4096 bytes
	V          [16]uint8      // array of registers => 1 register == 1 byte
	I          uint16         // index reg
	PC         uint16         // program counter
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

func (chip8 *Chip8) SetupInput() {

}

func (chip8 *Chip8) ClearDisplay() {
	for i := 0; i < 64*32; i++ {
		chip8.gfx[i] = 0
	}
}

func (chip8 *Chip8) Initialize() {
	chip8.PC = 0x200
	chip8.opcode = 0
	chip8.I = 0
	chip8.SP = 0

	for i := 0; i < 16; i++ {
		chip8.stack[i] = 0
		chip8.V[i] = 0
		chip8.key[i] = 0
	}

	chip8.ClearDisplay()

	for i := 0; i < 4096; i++ {
		chip8.memory[i] = 0
	}

	for i := 0; i < 80; i++ {
		chip8.memory[i] = fontset[i]
	}

	chip8.delayTimer = 0
	chip8.soundTimer = 0
	chip8.drawFlag = 1
}

func (chip8 *Chip8) LoadGame(gamePath string) {
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

func (chip8 *Chip8) EmulateCycle() {
	chip8.opcode = ((uint16(chip8.memory[chip8.PC]) << 8) | uint16(chip8.memory[chip8.PC+1]))

	switch chip8.opcode & 0xf000 {

	case 0x0000:
		switch chip8.opcode & 0x000f {
		case 0x0000:
			chip8.ClearDisplay()
			break

		case 0x000e:
			chip8.PC = chip8.stack[chip8.SP]
			chip8.SP--
			break

		default:
			fmt.Printf("Unknown opcode:0x%X\n", chip8.opcode)
			break
		}

	case 0x1000:
		chip8.PC = chip8.opcode & 0x0fff
		break

	case 0x2000:
		chip8.stack[chip8.SP] = chip8.PC
		chip8.SP++
		chip8.PC = chip8.opcode & 0x0fff
		break

	case 0x3000:
		if chip8.V[(chip8.opcode&0x0f00)>>8] == uint8(chip8.opcode&0x00ff) {
			chip8.PC += 2
		}
		break

	case 0x4000:
		if chip8.V[(chip8.opcode&0x0f00)>>8] != uint8(chip8.opcode&0x00ff) {
			chip8.PC += 2
		}
		break

	case 0x5000:
		if chip8.V[(chip8.opcode&0x0f00)>>8] == chip8.V[(chip8.opcode&0x00f0)>>4] {
			chip8.PC += 2
		}
		break

	case 0x6000:
		chip8.V[(chip8.opcode&0x0f00)>>8] = uint8(chip8.opcode & 0x00ff)
		break

	case 0x7000:
		chip8.V[(chip8.opcode&0x0f00)>>8] += uint8(chip8.opcode & 0x00ff)
		break

	case 0x8000:
		switch chip8.opcode & 0x000f {
		case 0x0000:
			chip8.V[(chip8.opcode&0x0f00)>>8] = chip8.V[(chip8.opcode&0x00f0)>>4]
			break

		case 0x0001:
			chip8.V[(chip8.opcode&0x0f00)>>8] |= chip8.V[(chip8.opcode&0x00f0)>>4]
			break

		case 0x0002:
			chip8.V[(chip8.opcode&0x0f00)>>8] &= chip8.V[(chip8.opcode&0x00f0)>>4]
			break

		case 0x0003:
			chip8.V[(chip8.opcode&0x0f00)>>8] ^= chip8.V[(chip8.opcode&0x00f0)>>4]
			break

		case 0x0004:
			chip8.V[(chip8.opcode&0x0f00)>>8] += chip8.V[(chip8.opcode&0x00f0)>>4]
			break

		case 0x0005:
			chip8.V[(chip8.opcode&0x0f00)>>8] -= chip8.V[(chip8.opcode&0x00f0)>>4]
			break

		case 0x0006:
			chip8.V[(chip8.opcode&0x0f00)>>8] = chip8.V[(chip8.opcode&0x0f00)>>8] / 2
			break

		case 0x0007:
			chip8.V[(chip8.opcode&0x0f00)>>8] = chip8.V[(chip8.opcode&0x00f0)>>4] - chip8.V[(chip8.opcode&0x0f00)>>8]
			break

		case 0x000e:
			chip8.V[(chip8.opcode&0x0f00)>>8] = chip8.V[(chip8.opcode&0x0f00)>>8] * 2
			break

		default:
			fmt.Printf("Unknown opcode:0x%X\n", chip8.opcode)
			break
		}

	case 0x0004:
		fmt.Println("case 0x0004")
		if chip8.V[(chip8.opcode&0x00f0)>>4] > (0xff - chip8.V[(chip8.opcode&0x0f00)>>8]) {
			chip8.V[0xf] = 1 // carry
		} else {
			chip8.V[0xf] = 0
		}

		chip8.V[(chip8.opcode&0x0f00)>>8] += chip8.V[(chip8.opcode&0x00f0)>>4]
		chip8.PC += 2
		break

	case 0x0033:
		fmt.Println("case 0x0033")
		chip8.memory[chip8.I] = chip8.V[(chip8.opcode&0x0f00)>>8] / 100
		chip8.memory[chip8.I+1] = (chip8.V[(chip8.opcode&0x0f00)>>8] / 10) % 10
		chip8.memory[chip8.I+2] = (chip8.V[(chip8.opcode&0x0f00)>>8] % 100) % 10
		chip8.PC += 2
		break

	case 0xa000:
		fmt.Println("case 0xa000")
		chip8.I = (chip8.opcode & 0xf000) >> 12
		chip8.PC += 2
		break

	case 0xd000:
		x := uint16(chip8.V[(chip8.opcode&0x0f00)>>8])
		y := uint16(chip8.V[(chip8.opcode&0x00f0)>>4])
		height := chip8.opcode & 0x000f
		var pixel uint8

		chip8.V[0xf] = 0

		for yline := uint16(0); yline < height; yline++ {
			pixel = chip8.memory[chip8.I+yline]

			for xline := uint16(0); xline < 8; xline++ {
				if pixel&(0x80>>xline) != 0 {
					if chip8.gfx[x+xline+((y+yline)*64)] == 1 {
						chip8.V[0xf] = 1
					}
					chip8.gfx[(x + xline + ((y + yline) * 64))] ^= 1
				}
			}
		}
		chip8.drawFlag = 1
		chip8.PC += 2
		break

	default:
		fmt.Printf("Unknown opcode:0x%X\n", chip8.opcode)
		break
	}

	if chip8.delayTimer > 0 {
		chip8.delayTimer--
	}

	if chip8.soundTimer > 0 {
		if chip8.soundTimer == 1 {
			fmt.Println("BEEP")
		}
		chip8.soundTimer--
	}

	chip8.PC += 2 // 1 opcode == 2 bytes => needs to increment by 2
}

func (chip8 *Chip8) SetKeys() {

}
