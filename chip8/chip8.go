package chip8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

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
	buf, err := ioutil.ReadFile(gamePath)
	util.CheckError(err)

	for i := 0; i < len(buf); i++ {
		chip8.memory[512+i] = buf[i]
	}
}

func (chip8 *Chip8) EmulateCycle() {
	chip8.opcode = ((uint16(chip8.memory[chip8.PC]) << 8) | uint16(chip8.memory[chip8.PC+1]))
	fmt.Printf("chip8.PC = 0x%X + opcode = 0x%X\n", chip8.PC, chip8.opcode)

	x := (chip8.opcode & 0x0f00) >> 8
	y := (chip8.opcode & 0x00f0) >> 4

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
		if chip8.V[x] == uint8(chip8.opcode&0x00ff) {
			chip8.PC += 2
		}
		break

	case 0x4000:
		if chip8.V[x] != uint8(chip8.opcode&0x00ff) {
			chip8.PC += 2
		}
		break

	case 0x5000:
		if chip8.V[x] == chip8.V[y] {
			chip8.PC += 2
		}
		break

	case 0x6000:
		chip8.V[x] = uint8(chip8.opcode & 0x00ff)
		break

	case 0x7000:
		chip8.V[x] += uint8(chip8.opcode & 0x00ff)
		break

	case 0x8000:
		switch chip8.opcode & 0x000f {
		case 0x0000:
			chip8.V[x] = chip8.V[y]
			break

		case 0x0001:
			chip8.V[x] |= chip8.V[y]
			break

		case 0x0002:
			chip8.V[x] &= chip8.V[y]
			break

		case 0x0003:
			chip8.V[x] ^= chip8.V[y]
			break

		case 0x0004:
			if chip8.V[y] > (0xff - chip8.V[x]) {
				chip8.V[0xf] = 1 // carry
			} else {
				chip8.V[0xf] = 0
			}

			chip8.V[x] += chip8.V[y]
			break

		case 0x0005:
			if chip8.V[y] > chip8.V[x] {
				chip8.V[0xf] = 1 // carry
			} else {
				chip8.V[0xf] = 0
			}

			chip8.V[x] -= chip8.V[y]
			break

		case 0x0006:
			if chip8.V[(chip8.opcode&0x0f00)]&0x1 == 1 {
				chip8.V[0xf] = 1
			} else {
				chip8.V[0xf] = 0
			}

			chip8.V[x] >>= 1
			break

		case 0x0007:
			if chip8.V[y] > chip8.V[x] {
				chip8.V[0xf] = 1
			} else {
				chip8.V[0xf] = 0
			}

			chip8.V[x] = chip8.V[y] - chip8.V[x]
			break

		case 0x000e:
			if chip8.V[(chip8.opcode&0x0f00)]&0x8 == 1 {
				chip8.V[0xf] = 1
			} else {
				chip8.V[0xf] = 0
			}

			chip8.V[x] <<= 1
			break

		default:
			fmt.Printf("Unknown opcode:0x%X\n", chip8.opcode)
			break
		}

	case 0x9000:
		if chip8.V[x] != chip8.V[y] {
			chip8.PC += 2
		}
		break

	case 0xa000:
		chip8.I = (chip8.opcode & 0x0fff)
		break

	case 0xb000:
		chip8.PC = uint16(chip8.V[0x0]) + (chip8.opcode & 0x0fff)
		break

	case 0xc000:
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		chip8.V[x] = uint8(r.Intn(255)) & chip8.V[(chip8.opcode&0x00ff)]
		break

	case 0xd000:
		xPos := uint16(chip8.V[x])
		yPos := uint16(chip8.V[y])
		height := chip8.opcode & 0x000f
		var pixel uint8

		chip8.V[0xf] = 0

		for yline := uint16(0); yline < height; yline++ {
			pixel = chip8.memory[chip8.I+yline]

			for xline := uint16(0); xline < 8; xline++ {
				if pixel&(0x80>>xline) != 0 {
					if chip8.gfx[xPos+xline+((yPos+yline)*64)] == 1 {
						chip8.V[0xf] = 1
					}
					chip8.gfx[(xPos + xline + ((yPos + yline) * 64))] ^= 1
				}
			}
		}
		chip8.drawFlag = 1
		break

	case 0xe000:
		switch chip8.opcode & 0x00ff {
		case 0x009e:
			if chip8.key[chip8.V[x]] != 0 {
				chip8.PC += 2
			}
			break

		case 0x00a1:
			if chip8.key[chip8.V[x]] == 0 {
				chip8.PC += 2
			}
			break
		}

	case 0xf000:
		switch chip8.opcode & 0x00ff {
		case 0x0007:
			chip8.V[x] = chip8.delayTimer
			break

		case 0x000a:
			break

		case 0x0015:
			chip8.delayTimer = chip8.V[x]
			break

		case 0x0018:
			chip8.soundTimer = chip8.V[x]
			break

		case 0x001e:
			chip8.I += uint16(chip8.V[x])
			break

		case 0x0029:
			chip8.I = uint16(chip8.V[x] * 5)
			break

		case 0x0033:
			chip8.memory[chip8.I] = chip8.V[x] / 100
			chip8.memory[chip8.I+1] = (chip8.V[x] / 10) % 10
			chip8.memory[chip8.I+2] = (chip8.V[x] % 100) % 10
			break

		case 0x0055:
			for i := uint16(0); i < x; i++ {
				chip8.memory[chip8.I+i] = chip8.V[i]
			}
			break

		case 0x0065:
			for i := uint16(0); i < x; i++ {
				chip8.V[i] = chip8.memory[chip8.I+i]
			}
			break
		}
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
