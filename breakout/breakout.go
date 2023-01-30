package main

import raylib "github.com/gen2brain/raylib-go/raylib"
import "math/rand"

const (
	BoardWidthInBricks  = 12
	BoardHeightInBricks = 13
	BrickWidthInPixels  = 64
	BrickHeightInPixels = 24
)

const (
	BrickOffsetX = 16
	BrickOffsetY = 16
)

const (
	None = iota - 1
	Left
	Top
	Right
	Bottom
)

type Brick struct {
	typeOf  int
	isAlive bool
}

type Rectangle struct {
	centerPosition raylib.Vector2
	size           raylib.Vector2
}

type Ball struct {
	Rectangle
	velocity raylib.Vector2
}

type InputScheme struct {
	leftButton  int32
	rightButton int32
}

type Pad struct {
	Rectangle
	InputScheme
	score    int
	velocity raylib.Vector2
}

var ball Ball
var player1 Pad
var bricks [BoardWidthInBricks][BoardHeightInBricks]*Brick

var InitialBallPosition raylib.Vector2
var InitialBallVelocity raylib.Vector2

func main() {
	raylib.InitWindow(800, 450, "GO Breakout")
	defer raylib.CloseWindow()
	raylib.SetTargetFPS(60)

	SetupGame()

	for !raylib.WindowShouldClose() {
		dt := raylib.GetFrameTime()
		Update(dt)
		Draw()
	}
}

func SetupGame() {
	screenSizeX := raylib.GetScreenWidth()
	screenSizeY := raylib.GetScreenHeight()

	{ // Setup bricks
		for i := 0; i < BoardWidthInBricks; i++ {
			for j := 0; j < BoardHeightInBricks; j++ {
				bricks[i][j] = new(Brick)
				bricks[i][j].typeOf = rand.Intn(4)
				bricks[i][j].isAlive = true
			}
		}
	}
	{ // Set up ball
		InitialBallPosition = raylib.Vector2{float32(screenSizeX / 2), float32(screenSizeY - 20)}
		InitialBallVelocity = raylib.Vector2{50, -25}
		ball.velocity = InitialBallVelocity
		ball.centerPosition = InitialBallPosition
		ball.size = raylib.Vector2{10, 10}
	}
	{ // Set up player
		player1.size = raylib.Vector2{50, 5}
		player1.velocity = raylib.Vector2{100, 100}
		player1.centerPosition = raylib.Vector2{float32(screenSizeX / 2), float32(screenSizeY - 10)}
		player1.InputScheme = InputScheme{
			raylib.KeyA,
			raylib.KeyD,
		}
	}
}

func Update(deltaTime float32) {
	height := raylib.GetScreenHeight()
	width := raylib.GetScreenWidth()
	collisionFace := None

	{ // Update Player
		if raylib.IsKeyDown(player1.rightButton) {
			// Update position
			player1.centerPosition.X += (deltaTime * player1.velocity.X)
			// Clamp on right edge
			if player1.centerPosition.X+(player1.size.X/2) > float32(width) {
				player1.centerPosition.X = float32(width) - (player1.size.X / 2)
			}
		}
		if raylib.IsKeyDown(player1.leftButton) {
			// Update position
			player1.centerPosition.X -= (deltaTime * player1.velocity.X)
			// Clamp on left edge
			if player1.centerPosition.X-(player1.size.X/2) < 0 {
				player1.centerPosition.X = (player1.size.X / 2)
			}
		}
	}
	{ // Update ball
		ball.centerPosition.X += deltaTime * ball.velocity.X
		ball.centerPosition.Y += deltaTime * ball.velocity.Y
	}
	// Collisions
	{ // ball boundary collisions
		isBallOnBottomScreenEdge := ball.centerPosition.Y > float32(height)
		isBallOnTopScreenEdge := ball.centerPosition.Y < float32(0)
		isBallOnLeftRightScreenEdge := ball.centerPosition.X > float32(width) || ball.centerPosition.X < float32(0)
		if isBallOnBottomScreenEdge {
			ball.centerPosition = InitialBallPosition
			ball.velocity = InitialBallVelocity
		}
		if isBallOnTopScreenEdge {
			ball.velocity.Y *= -1
		}
		if isBallOnLeftRightScreenEdge {
			ball.velocity.X *= -1
		}
	}
	{ // ball brick collisions
		hasHit := false
		for i := 0; i < BoardWidthInBricks; i++ {
			for j := 0; j < BoardHeightInBricks; j++ {
				brick := bricks[i][j]
				if !brick.isAlive {
					continue
				}

				// Coords
				brickX := float32(BrickOffsetX + (i * BrickWidthInPixels))
				brickY := float32(BrickOffsetY + (j * BrickHeightInPixels))

				// Ball position
				ballX := ball.centerPosition.X - (ball.size.X / 2)
				ballY := ball.centerPosition.Y - (ball.size.Y / 2)

				// Center Brick
				brickCenterX := brickX + (BrickWidthInPixels / 2)
				brickCenterY := brickY + (BrickHeightInPixels / 2)

				hasCollisionX := ballX+ball.size.X >= brickX && brickX+BrickWidthInPixels >= ballX
				hasCollisionY := ballY+ball.size.Y >= brickY && brickY+BrickHeightInPixels >= ballY

				if hasCollisionX && hasCollisionY {
					brick.isAlive = false
					hasHit = true

					// Determine which face of the brick was hit
					ymin := Max(brickY, ballY)
					ymax := Min(brickY+BrickHeightInPixels, ballY+ball.size.Y)
					ysize := ymax - ymin
					xmin := Max(brickX, ballX)
					xmax := Min(brickX+BrickWidthInPixels, ballX+ball.size.X)
					xsize := xmax - xmin
					if xsize > ysize && ball.centerPosition.Y > brickCenterY {
						collisionFace = Bottom
					} else if xsize > ysize && ball.centerPosition.Y <= brickCenterY {
						collisionFace = Top
					} else if xsize <= ysize && ball.centerPosition.X > brickCenterX {
						collisionFace = Right
					} else if xsize <= ysize && ball.centerPosition.X <= brickCenterX {
						collisionFace = Left
					} else {
						// Could assert or panic here
					}

					break
				}
			}
			if hasHit {
				break
			}
		}
	}
	{ // Update ball after collision
		if collisionFace != None {
			hasPositiveX := ball.velocity.X > 0
			hasPositiveY := ball.velocity.Y > 0
			if (collisionFace == Top && hasPositiveX && hasPositiveY) ||
				(collisionFace == Top && !hasPositiveX && hasPositiveY) ||
				(collisionFace == Bottom && hasPositiveX && !hasPositiveY) ||
				(collisionFace == Bottom && !hasPositiveX && !hasPositiveY) {
				ball.velocity.Y *= -1
			}
			if (collisionFace == Left && hasPositiveX && hasPositiveY) ||
				(collisionFace == Left && hasPositiveX && !hasPositiveY) ||
				(collisionFace == Right && !hasPositiveX && hasPositiveY) ||
				(collisionFace == Right && !hasPositiveX && !hasPositiveY) {
				ball.velocity.X *= -1
			}
		}
	}
	{ // Update ball after pad collision
		if DetectBallTouchesPad(ball, &player1) {
			previousVelocity := ball.velocity
			distanceX := ball.centerPosition.X - player1.centerPosition.X
			percentage := distanceX / (player1.size.X / 2)
			ball.velocity.X = InitialBallVelocity.X * percentage
			ball.velocity.Y *= -1
			newVelocity := raylib.Vector2Scale(raylib.Vector2Normalize(ball.velocity), (raylib.Vector2Length(previousVelocity) * 1.1))
			ball.velocity = newVelocity
		}
	}
	{ // Detect all bricks popped
		hasAtLeastOneBrick := false
		for i := 0; i < BoardWidthInBricks; i++ {
			for j := 0; j < BoardHeightInBricks; j++ {
				brick := bricks[i][j]
				if brick.isAlive {
					hasAtLeastOneBrick = true
					break	// NOTE: This needs to break all the way out to be a proper comparison of identical code execution
				}
			}
		}
		if !hasAtLeastOneBrick {
			SetupGame()
		}
	}
}

func Draw() {
	raylib.BeginDrawing()
	defer raylib.EndDrawing()
	raylib.ClearBackground(raylib.Black)

	{ // Draw alive bricks
		for i := 0; i < BoardWidthInBricks; i++ {
			for j := 0; j < BoardHeightInBricks; j++ {
				if !bricks[i][j].isAlive {
					continue
				}

				raylib.DrawRectangle(int32(BrickOffsetX+(i*BrickWidthInPixels)), int32(BrickOffsetY+(j*BrickHeightInPixels)), BrickWidthInPixels, BrickHeightInPixels, TypeToColor(bricks[i][j].typeOf))
			}
		}
	}
	{ // Draw Players
		raylib.DrawRectangle(int32(player1.centerPosition.X-(player1.size.X/2)), int32(player1.centerPosition.Y-(player1.size.Y/2)), int32(player1.size.X), int32(player1.size.Y), raylib.White)
	}
	{ // Draw Ball
		raylib.DrawRectangle(int32(ball.centerPosition.X-(ball.size.X/2)), int32(ball.centerPosition.Y-(ball.size.Y/2)), int32(ball.size.X), int32(ball.size.Y), raylib.White)
	}
}

func DetectBallTouchesPad(ball Ball, pad *Pad) bool {
	ballX := ball.centerPosition.X - (ball.size.X / 2)
	ballY := ball.centerPosition.Y - (ball.size.Y / 2)
	padX := pad.centerPosition.X - (pad.size.X / 2)
	padY := pad.centerPosition.Y - (pad.size.Y / 2)
	if ballY+(ball.size.Y/2) >= padY && ballX >= padX && ballX <= padX+pad.size.X {
		return true
	}
	return false
}

func TypeToColor(typeOf int) raylib.Color {
	switch typeOf {
	case 0:
		return raylib.White
	case 1:
		return raylib.Red
	case 2:
		return raylib.Green
	case 3:
		return raylib.Blue
	}
	return raylib.Color{}
}

func Max(a float32, b float32) float32 { // Yes Math really doesn't have a max for float32.
	if a > b {
		return a
	}
	return b
}

func Min(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
