package internal

import (
	"fmt"
	"io/ioutil"
	"math/rand"
)

var (
	DebugMode = true
)

//CHIP8 emulator
type CHIP8 struct {
	//current 16 bit opecode
	opcode uint16
	//memory from 0 to 4096
	memory [4096]uint8
	//16 8-bit registers
	v [16]uint8
	//16 nesting level
	stack [16]uint16
	sp    uint16
	//program counter range from 0 to 4096(0xFFF)
	pc uint16
	//register index range from 0 to 4096(0xFFF)
	i uint16
	//screen
	Gfx [64 * 32]uint8
	//timer
	delayTimer uint8
	soundTimer uint8
	//hex keys
	key [16]uint8
	//draw flag
	drawFlag bool
}

func debugPrintln(v ...interface{}) {
	if DebugMode {
		fmt.Println(v...)
	}
}
func debugPrintf(format string, v ...interface{}) {
	if DebugMode {
		fmt.Printf(format, v...)
	}
}

func (c *CHIP8) Initialize() error {
	//pc start at 0x200
	c.pc = 0x200

	//clear display
	//clear stack
	//clear register
	//clear memory

	//load fontset to memory from 0 to 80
	var chip8Fontset = []byte{
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

	for i := 0; i < len(chip8Fontset); i++ {
		c.memory[i] = chip8Fontset[i]
	}

	return nil
}
func (c *CHIP8) LoadRom(romDir string) error {
	romByte, err := ioutil.ReadFile(romDir)
	if err != nil {
		return err
	}
	//load to memory
	for i := 0; i < len(romByte); i++ {
		c.memory[512+i] = romByte[i]
	}
	debugPrintln(c.memory[512:])
	return nil
}
func (c *CHIP8) EmulateCycle() {
	//fetch code
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	//decode code
	//excute code
	debugPrintf("0x%X\n", c.opcode)
	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode {
		case 0x0E0:
			debugPrintln("display Clears the screen.")
			c.Gfx = [64 * 32]uint8{}
			c.pc += 2
		case 0x00EE:
			debugPrintln("flow return return;	Returns from a subroutine.")
			c.pc = c.stack[c.sp] + 2
			c.sp--
		default:
			debugPrintln("unknow call")
			c.pc += 2
		}

	case 0x1000:
		debugPrintln("Flow goto NNN;")
		c.pc = c.opcode & 0x0FFF
	case 0x2000:
		debugPrintln("Flow *(0xNNN)()")
		c.sp++
		c.stack[c.sp] = c.pc
		c.pc = c.opcode & 0x0FFF
	case 0x3000:
		debugPrintln("Cond	if(Vx==NN)")
		x := (c.opcode & 0x0F00) >> 8
		nn := uint8(c.opcode & 0x00FF)
		if c.v[x] == nn {
			c.pc += 4
		} else {
			c.pc += 2
		}

	case 0x4000:
		debugPrintln("Cond	if(Vx!=NN)")
		x := (c.opcode & 0x0F00) >> 8
		nn := uint8(c.opcode & 0x00FF)
		if c.v[x] != nn {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x5000:
		debugPrintln("Cond	if(Vx==Vy)")
		x := (c.opcode & 0x0F00) >> 8
		y := (c.opcode & 0x00F0) >> 4

		if c.v[x] == c.v[y] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x6000:
		debugPrintln("Const	Vx = NN")
		x := (c.opcode & 0x0F00) >> 8
		nn := uint8(c.opcode & 0x00FF)
		c.v[x] = nn
		c.pc += 2
	case 0x7000:
		debugPrintln("Const	Vx += NN")
		x := (c.opcode & 0x0F00) >> 8
		nn := uint8(c.opcode & 0x00FF)
		c.v[x] += nn
		c.pc += 2
	case 0x8000:
		x := (c.opcode & 0x0F00) >> 8
		y := (c.opcode & 0x00F0) >> 4

		switch c.opcode & 0x000F {
		case 0x0:
			debugPrintln("Assign	Vx=Vy")
			c.v[x] = c.v[y]
			c.pc += 2
		case 0x1:
			debugPrintln("BitOp	Vx=Vx|Vy")
			c.v[x] = c.v[x] | c.v[y]
			c.pc += 2
		case 0x2:
			debugPrintln("BitOp	Vx=Vx&Vy")
			c.v[x] = c.v[x] & c.v[y]
			c.pc += 2
		case 0x3:
			debugPrintln("BitOp	Vx=Vx^Vy")
			c.v[x] = c.v[x] ^ c.v[y]
			c.pc += 2
		case 0x4:
			debugPrintln("Math	Vx += Vy")
			if c.v[x]+c.v[y] > 0xFF {
				c.v[0xF] = 1 //carray
			} else {
				c.v[0xF] = 0
			}
			c.v[x] += c.v[y]
			c.pc += 2
		case 0x5:
			debugPrintln("Math	Vx -= Vy")
			if c.v[x] < c.v[y] {
				c.v[0xF] = 1 //borrow
			} else {
				c.v[0xF] = 0
			}
			c.v[x] -= c.v[y]
			c.pc += 2
		case 0x6:
			fmt.Println("bitOP vx>>=1")
			debugPrintln("BitOp	Vx>>=1")
			c.v[0xF] = c.v[x] & 0x0F
			c.v[x] = c.v[x] >> 1
			c.pc += 2
		case 0x7:
			debugPrintln("Math	Vx=Vy-Vx")
			if c.v[x] > c.v[y] {
				c.v[0xF] = 0 //borrow
			} else {
				c.v[0xF] = 1
			}
			c.v[x] = c.v[y] - c.v[x]
			c.pc += 2
		case 0xE:
			debugPrintln("BitOp	Vx<<=1")
			c.v[0XF] = (c.v[x] & 0xF0) >> 1
			c.v[x] = c.v[x] << 1
			c.pc += 2
		default:
			debugPrintf("unknow opcode: 0x%x\n", c.opcode)
			c.pc += 2
		}
	case 0x9000:
		debugPrintln("Cond	if(Vx!=Vy)")
		x := (c.opcode & 0x0F00) >> 8
		y := (c.opcode & 0x00F0) >> 4
		if c.v[x] != c.v[y] {
			c.pc += 4
		} else {
			c.pc += 2
		}

	case 0xA000:
		debugPrintln("MEM	I = NNN")
		c.i = c.opcode & 0x0FFF
		c.pc += 2
	case 0xB000:
		debugPrintln("Flow	PC=V0+NNN")
		c.pc = uint16(c.v[0]) + c.opcode&0x0FFF
	case 0xC000:
		debugPrintln("Rand	Vx=rand()&NN")
		x := (c.opcode & 0x0F00) >> 8
		c.v[x] = uint8(rand.Intn(255)) & uint8(c.opcode&0x0FF)
		c.pc += 2
	case 0xD000:

		x := uint16(c.v[(c.opcode&0x0F00)>>8])
		y := uint16(c.v[(c.opcode&0x00F0)>>4])
		height := c.opcode & 0x000F
		debugPrintln("Disp	draw(Vx,Vy,N)", x, y, height)
		c.v[0xF] = 0
		var pixels uint8
		for row := uint16(0); row < height; row++ {
			pixels = c.memory[c.i+row]
			//check each pixel
			for col := uint16(0); col < 8; col++ {
				index := (row+y)*64 + col + x
				if index > 2048 {
					continue
				}
				//pixel not empty
				//0x80 represents as binary is 10000000
				if pixels&(0x80>>col) != 0 {
					//chech col if overlap
					if (row+y) > 32 || (col+x) > 64 {
						continue
					}

					if c.Gfx[index] == 1 {
						c.v[0xF] = 1
					}
					c.Gfx[index] ^= 1
				}
			}
		}
		//fmt.Println(c.Gfx)
		c.drawFlag = true
		c.pc += 2
	case 0xE000:
		switch c.opcode & 0x00FF {
		case 0x9E:
			debugPrintln("KeyOp	if(key()==Vx)")
			fmt.Println("KeyOp	if(key()==Vx)")
			x := (c.opcode & 0x0F00) >> 8
			if c.key[c.v[x]] != 0 {
				//does need to unset keymap?
				//c.key[c.v[x]]=0
				c.pc += 4
			} else {
				c.pc += 2
			}

		case 0xA1:
			debugPrintln("KeyOp	if(key()!=Vx)")
			x := (c.opcode & 0x0F00) >> 8
			//fmt.Println("KeyOp	if(key()!=Vx)", c.v[x])
			if c.key[c.v[x]] == 0 {
				c.pc += 4
			} else {
				c.pc += 2
				c.key[c.v[x]] = 0
			}
		default:
			debugPrintf("unknow opcode: 0x%x\n", c.opcode)
			c.pc += 2
		}
	case 0xF000:
		switch c.opcode & 0x00FF {
		case 0x07:
			debugPrintln("Timer	Vx = get_delay()")
			x := (c.opcode & 0x0F00) >> 8
			c.v[x] = c.delayTimer
			c.pc += 2
		case 0x0A:
			fmt.Println("get key")
			debugPrintln("KeyOp	Vx = get_key()")
			x := (c.opcode & 0x0F00) >> 8
			for index, k := range c.key {
				if k != 0 {
					c.v[x] = byte(index)
					c.pc += 2
					break
				}
			}
			c.key[c.v[x]] = 0
			c.pc += 2
		case 0x15:
			debugPrintln("Timer	delay_timer(Vx)")
			x := (c.opcode & 0x0F00) >> 8
			c.delayTimer = c.v[x]
			c.pc += 2
		case 0x18:
			debugPrintln("Sound	sound_timer(Vx)")
			x := (c.opcode & 0x0F00) >> 8
			c.soundTimer = c.v[x]
			c.pc += 2
		case 0x1E:
			debugPrintln("MEM	I +=Vx")
			x := (c.opcode & 0x0F00) >> 8
			c.i += uint16(c.v[x])
			c.pc += 2
		case 0x29:
			debugPrintln("MEM	I=sprite_addr[Vx]")
			x := (c.opcode & 0x0F00) >> 8
			//??Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			c.i = uint16(c.v[x] * 5)
			c.pc += 2
		case 0x33:
			debugPrintln("BCD")
			x := (c.opcode & 0x0F00) >> 8
			c.memory[c.i] = c.v[x] / 100
			c.memory[c.i+1] = (c.v[x] / 10) % 10
			c.memory[c.i+2] = (c.v[x] % 100) % 10
			c.pc += 2
		case 0x55:
			debugPrintln("MEM	reg_dump(Vx,&I)")
			x := (c.opcode & 0x0F00) >> 8
			for j := uint16(0); j <= x; j++ {
				c.memory[c.i+j] = c.v[j]
			}
			c.pc += 2
		case 0x65:
			debugPrintln("MEM	reg_load(Vx,&I)")
			x := (c.opcode & 0x0F00) >> 8
			for j := uint16(0); j <= x; j++ {
				c.v[j] = c.memory[c.i+j]
			}
			c.pc += 2
		default:
			debugPrintf("unknow opcode: 0x%x\n", c.opcode)
			c.pc += 2
		}
	default:
		debugPrintf("unknow opcode: 0x%x\n", c.opcode)
		c.pc += 2
	}

	//set timer
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.soundTimer--
	}
}

func (c *CHIP8) DrawFlag() bool {
	flag := c.drawFlag
	if flag {
		c.drawFlag = false
	}
	return flag

}

func (c *CHIP8) SetKeys(i int) {
	c.key[i] = 1
}
