package main

import (
	"chip8/internal"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	squareSize = 1
	cols       = 64
	rows       = 32
	width      = cols * squareSize
	height     = rows * squareSize
)

func drawCell(screen [2048]uint8) {

	for y := int32(0); y < rows; y++ {
		linePos := y * squareSize
		for x := int32(0); x < cols; x++ {
			cellV := screen[y*64+x]

			if cellV == 1 {
				rl.DrawRectangle(x*squareSize, linePos, squareSize, squareSize, rl.LightGray)
			}

		}
	}
}

func main() {
	fmt.Println("start...")
	internal.DebugMode = false
	emu := &internal.CHIP8{}
	err := emu.Initialize()
	if err != nil {
		panic(err)
	}
	err = emu.LoadRom("./pong.ch8")
	if err != nil {
		panic(err)
	}

	rl.InitWindow(width, height, "CHIP-8")
	rl.SetTargetFPS(60)

	//main loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		emu.EmulateCycle()

		if emu.DrawFlag() {
			//draw
			rl.ClearBackground(rl.Black)
			drawCell(emu.Gfx)
		}
		if rl.IsKeyDown(rl.KeyW) {
			emu.SetKeys(1)
		}
		if rl.IsKeyDown(rl.KeyS) {
			emu.SetKeys(4)
		}
		if rl.IsKeyDown(rl.KeyUp) {
			emu.SetKeys(12)
		}
		if rl.IsKeyDown(rl.KeyDown) {
			emu.SetKeys(13)
		}
		//emu.SetKeys(1)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
