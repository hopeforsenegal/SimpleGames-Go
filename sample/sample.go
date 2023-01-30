package main

import "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(800, 450, "GO Sample")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.White)

		rl.EndDrawing()
	}
}