package main

import raylib "github.com/gen2brain/raylib-go/raylib"
import "strconv"

type TextAlignment int64

const (
	Left TextAlignment = iota
	Center
	Right
)

type Rectangle struct {
	centerPosition raylib.Vector2
	size           raylib.Vector2
}

type Ball struct {
	Rectangle
	velocity raylib.Vector2
}

type InputScheme struct {
	upButton   int32
	downButton int32
}

type Pad struct {
	Rectangle
	InputScheme
	score    int
	velocity raylib.Vector2
}

var ball Ball
var player1 Pad
var player2 Pad
var players []*Pad = []*Pad{&player1, &player2}

var InitialBallPosition raylib.Vector2

func main() {
	raylib.InitWindow(800, 450, "GO Pong")
	defer raylib.CloseWindow()
	raylib.SetTargetFPS(60)

	screenSizeX := raylib.GetScreenWidth()
	screenSizeY := raylib.GetScreenHeight()

	InitialBallPosition = raylib.Vector2{float32(screenSizeX / 2), float32(screenSizeY / 2)}
	ball.velocity = raylib.Vector2{50, 25}
	ball.centerPosition = InitialBallPosition
	ball.size = raylib.Vector2{10, 10}
	player2.size = raylib.Vector2{5, 50}
	player1.size = raylib.Vector2{5, 50}
	player2.velocity = raylib.Vector2{100, 100}
	player1.velocity = raylib.Vector2{100, 100}
	player1.centerPosition = raylib.Vector2{float32(0 + 5), float32(screenSizeY / 2)}
	player2.centerPosition = raylib.Vector2{float32(float32(screenSizeX) - player2.size.X - 5), float32(screenSizeY / 2)}
	player1.InputScheme = InputScheme{
		raylib.KeyW,
		raylib.KeyS,
	}
	player2.InputScheme = InputScheme{
		raylib.KeyI,
		raylib.KeyK,
	}

	for !raylib.WindowShouldClose() {
		dt := raylib.GetFrameTime()
		Update(dt)
		Draw()
	}
}

func Update(deltaTime float32) {
	height := raylib.GetScreenHeight()
	width := raylib.GetScreenWidth()
	{ // Update players
		for _, player := range players {
			if raylib.IsKeyDown(player.downButton) {
				// Update position
				player.centerPosition.Y += (deltaTime * player.velocity.Y)
				// Clamp on bottom edge
				if player.centerPosition.Y+(player.size.Y/2) > float32(height) {
					player.centerPosition.Y = float32(height) - (player.size.Y / 2)
				}
			}
			if raylib.IsKeyDown(player.upButton) {
				// Update position
				player.centerPosition.Y -= (deltaTime * player.velocity.Y)
				// Clamp on top edge
				if player.centerPosition.Y-(player.size.Y/2) < 0 {
					player.centerPosition.Y = (player.size.Y / 2)
				}
			}
		}
	}
	{ // Update ball
		ball.centerPosition.X += deltaTime * ball.velocity.X
		ball.centerPosition.Y += deltaTime * ball.velocity.Y
	}
	{ // Check collisions
		for _, player := range players {
			isDetectBallTouchesPad := DetectBallTouchesPad(ball, player)
			if isDetectBallTouchesPad {
				ball.velocity.X *= -1
			}
		}
		isBallOnTopBottomScreenEdge := ball.centerPosition.Y > float32(height) || ball.centerPosition.Y < 0
		isBallOnRightScreenEdge := ball.centerPosition.X > float32(width)
		isBallOnLeftScreenEdge := ball.centerPosition.X < 0
		if isBallOnTopBottomScreenEdge {
			ball.velocity.Y *= -1
		}
		if isBallOnLeftScreenEdge {
			ball.centerPosition = InitialBallPosition
			player2.score += 1
		}
		if isBallOnRightScreenEdge {
			ball.centerPosition = InitialBallPosition
			player1.score += 1
		}
	}
}

func Draw() {
	raylib.BeginDrawing()
	defer raylib.EndDrawing()
	raylib.ClearBackground(raylib.Black)

	{ // Draw Court Line
		var LineThinkness float32 = 2.0
		x := float32(raylib.GetScreenWidth() / 2.0)
		from := raylib.Vector2{x, 5.0}
		to := raylib.Vector2{x, float32(raylib.GetScreenHeight() - 5.0)}
		raylib.DrawLineEx(from, to, LineThinkness, raylib.LightGray)
	}
	{ // Draw Scores
		DrawText(strconv.Itoa(player1.score), Right, int32(raylib.GetScreenWidth()/2)-10, 10, 20)
		DrawText(strconv.Itoa(player2.score), Left, int32(raylib.GetScreenWidth()/2)+10, 10, 20)
	}
	{ // Draw Players
		for _, player := range players {
			raylib.DrawRectangle(int32(player.centerPosition.X-(player.size.X/2)), int32(player.centerPosition.Y-(player.size.Y/2)), int32(player.size.X), int32(player.size.Y), raylib.White)
		}
	}
	{ // Draw Ball
		raylib.DrawRectangle(int32(ball.centerPosition.X-(ball.size.X/2)), int32(ball.centerPosition.Y-(ball.size.Y/2)), int32(ball.size.X), int32(ball.size.Y), raylib.White)
	}
}

func DetectBallTouchesPad(ball Ball, pad *Pad) bool {
	if ball.centerPosition.X >= pad.centerPosition.X && ball.centerPosition.X <= pad.centerPosition.X+pad.size.X {
		if ball.centerPosition.Y >= pad.centerPosition.Y-(pad.size.Y/2) && ball.centerPosition.Y <= pad.centerPosition.Y+pad.size.Y/2 {
			return true
		}
	}
	return false
}

func DrawText(text string, alignment TextAlignment, posX int32, posY int32, fontSize int32) {
	fontColor := raylib.LightGray
	if alignment == Left {
		raylib.DrawText(text, posX, posY, fontSize, fontColor)
	} else if alignment == Center {
		scoreSizeLeft := raylib.MeasureText(text, fontSize)
		raylib.DrawText(text, posX-(scoreSizeLeft/2), posY, fontSize, fontColor)
	} else if alignment == Right {
		scoreSizeLeft := raylib.MeasureText(text, fontSize)
		raylib.DrawText(text, posX-scoreSizeLeft, posY, fontSize, fontColor)
	}
}
