package chip8

// Chip8 structure
type Chip8 struct {
	opcode [2]byte    // 1 opcode == 2 bytes
	memory [4096]byte // memory == 4096 bytes
	V      [8]byte    // array of registers => 1 register == 1 byte
	I      [2]byte
	PC     [2]byte
}
