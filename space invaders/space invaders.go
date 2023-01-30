package main

import raylib "github.com/gen2brain/raylib-go/raylib"
import "math/rand"
import "strconv"

type TextAlignment int64

const (
	Left TextAlignment = iota
	Center
	Right
)

const (
	BulletCooldownSeconds = 0.3
	MaxNumBullets         = 50
	MaxNumEnemies         = 50
)

type Rectangle struct {
	centerPosition raylib.Vector2
	size           raylib.Vector2
}

type InputScheme struct {
	leftButton  int32
	rightButton int32
	shootButton int32
}

type Pad struct {
	Rectangle
	InputScheme
	score    int
	velocity raylib.Vector2
}

type Bullet struct {
	Rectangle
	velocity raylib.Vector2
	isActive bool
	color    raylib.Color
}
type Enemy struct {
	Rectangle
	velocity raylib.Vector2
	isActive bool
	color    raylib.Color
}

var bullets [MaxNumBullets]*Bullet
var enemies [MaxNumEnemies]*Enemy
var player1 Pad
var m_TimerBulletCooldown float32
var m_TimerSpawnEnemy float32
var numEnemiesThisLevel int
var numEnemiesToSpawn int
var numEnemiesKilled int
var numLives = 3
var IsGameOver bool
var IsWin bool

var InitialPlayerPosition raylib.Vector2

func main() {
	raylib.InitWindow(800, 450, "GO Space Invaders")
	defer raylib.CloseWindow()
	raylib.SetTargetFPS(60)

	screenSizeX := raylib.GetScreenWidth()
	screenSizeY := raylib.GetScreenHeight()
	InitialPlayerPosition = raylib.Vector2{float32(screenSizeX / 2), float32(screenSizeY - 10)}

	{ // Set up player
		player1.size = raylib.Vector2{25, 25}
		player1.velocity = raylib.Vector2{100, 100}
		player1.centerPosition = InitialPlayerPosition
		player1.InputScheme = InputScheme{
			raylib.KeyA,
			raylib.KeyD,
			raylib.KeySpace,
		}
	}
	{ // init bullets
		for i := 0; i < MaxNumBullets; i++ {
			bullets[i] = new(Bullet)
			{
				bullets[i].velocity = raylib.Vector2{0, 400}
				bullets[i].Rectangle = Rectangle{size: raylib.Vector2{5, 5}}
			}
		}
	}
	{ // init enemies
		for i := 0; i < MaxNumEnemies; i++ {
			enemies[i] = new(Enemy)
			{
				enemies[i].velocity = raylib.Vector2{0, 40}
				enemies[i].Rectangle = Rectangle{
					raylib.Vector2{float32(rand.Intn(screenSizeX)), -20},
					raylib.Vector2{20, 20},
				}
			}
		}
		numEnemiesToSpawn = 10
		numEnemiesThisLevel = 10
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

	if IsGameOver || IsWin {
		return
	}

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
		if HasHitTime(&m_TimerBulletCooldown, deltaTime) {
			if raylib.IsKeyDown(player1.shootButton) {
				for i := 0; i < MaxNumBullets; i++ {
					if !bullets[i].isActive {
						m_TimerBulletCooldown = BulletCooldownSeconds
						bullets[i].isActive = true
						{
							bullets[i].centerPosition.X = player1.centerPosition.X
							bullets[i].centerPosition.Y = player1.centerPosition.Y + (player1.size.Y / 4)
							break
						}
					}
				}
			}
		}
	}
	{ // Update active bullets
		for i := 0; i < MaxNumBullets; i++ {
			bullet := bullets[i]
			// Movement
			if bullet.isActive {
				bullet.centerPosition.Y -= bullet.velocity.Y * deltaTime

				// Went off screen
				if bullet.centerPosition.Y+(bullet.size.Y/2) <= 0 {
					bullet.isActive = false
				}
			}
		}
	}
	{ // Update active enemies
		for i := 0; i < numEnemiesThisLevel; i++ {
			enemy := enemies[i]
			// Movement
			if enemy.isActive {
				enemy.centerPosition.Y += enemy.velocity.Y * deltaTime

				// Went off screen
				if enemy.centerPosition.Y-(enemy.size.Y/2) >= float32(height) {
					enemy.centerPosition = raylib.Vector2{float32(rand.Intn(width)), -20}
				} else {
					enemyX := enemy.centerPosition.X - (enemy.size.X / 2)
					enemyY := enemy.centerPosition.Y - (enemy.size.Y / 2)
					{ // bullet | enemy collision
						for j := 0; j < MaxNumBullets; j++ {
							bullet := bullets[j]
							bulletX := bullet.centerPosition.X - (bullet.size.X / 2)
							bulletY := bullet.centerPosition.Y - (bullet.size.Y / 2)

							hasCollisionX := bulletX+bullet.size.X >= enemyX && enemyX+enemy.size.X >= bulletX
							hasCollisionY := bulletY+bullet.size.Y >= enemyY && enemyY+enemy.size.Y >= bulletY

							if hasCollisionX && hasCollisionY {
								bullet.isActive = false
								enemy.isActive = false
								{
									numEnemiesKilled++
									IsWin = numEnemiesKilled >= numEnemiesThisLevel
									break
								}
							}
						}
					}
					{ // player | enemy collision
						bulletX := player1.centerPosition.X - (player1.size.X / 2)
						bulletY := player1.centerPosition.Y - (player1.size.Y / 2)

						hasCollisionX := bulletX+player1.size.X >= enemyX && enemyX+enemy.size.X >= bulletX
						hasCollisionY := bulletY+player1.size.Y >= enemyY && enemyY+enemy.size.Y >= bulletY

						if hasCollisionX && hasCollisionY {
							enemy.isActive = false
							{
								player1.centerPosition = InitialPlayerPosition
								numLives--
								IsGameOver = numLives <= 0
							}
						}
					}
				}
			}
		}
	}
	{ // Spawn enemies
		canSpawn := HasHitInterval(&m_TimerSpawnEnemy, 2.0, deltaTime)
		for i := 0; i < MaxNumEnemies; i++ {
			enemy := enemies[i]
			// Spawn
			if !enemy.isActive {
				if canSpawn && numEnemiesToSpawn > 0 {
					numEnemiesToSpawn--
					enemy.isActive = true
					{
						enemy.centerPosition = raylib.Vector2{float32(rand.Intn(width)), -20}
						break
					}
				}
			}
		}
	}
}

func Draw() {
	raylib.BeginDrawing()
	defer raylib.EndDrawing()
	raylib.ClearBackground(raylib.White)

	height := int32(raylib.GetScreenHeight())
	width := int32(raylib.GetScreenWidth())

	{ // Draw Players
		raylib.DrawRectangle(int32(player1.centerPosition.X-(player1.size.X/2)), int32(player1.centerPosition.Y-(player1.size.Y/2)), int32(player1.size.X), int32(player1.size.Y), raylib.Black)
	}
	{ // Draw the bullets
		for i := 0; i < MaxNumBullets; i++ {
			bullet := bullets[i]
			if bullet.isActive {
				raylib.DrawRectangle(int32(bullet.centerPosition.X-(bullet.size.X/2)),
					int32(bullet.centerPosition.Y-(bullet.size.Y/2)),
					int32(bullet.size.X),
					int32(bullet.size.Y),
					raylib.Orange)
			}
		}
	}
	{ // Draw the enemies
		for i := 0; i < MaxNumEnemies; i++ {
			enemy := enemies[i]
			if enemy.isActive {
				raylib.DrawRectangle(int32(enemy.centerPosition.X-(enemy.size.X/2)),
					int32(enemy.centerPosition.Y-(enemy.size.Y/2)),
					int32(enemy.size.X),
					int32(enemy.size.Y),
					raylib.Blue)
			}
		}
	}
	{ // Draw Info
		DrawText("Lives "+strconv.Itoa(numLives), Left, 15, 5, 20)

		if IsGameOver {
			DrawText("Game Over", Center, width/2, height/2, 50)
		}
		if IsWin {
			DrawText("You Won", Center, width/2, height/2, 50)
		}
	}
}

func DrawText(text string, alignment TextAlignment, posX int32, posY int32, fontSize int32) {
	fontColor := raylib.DarkGray
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

func HasHitInterval(timeRemaining *float32, resetTime float32, deltaTime float32) bool {
	*timeRemaining -= deltaTime
	if *timeRemaining <= 0 {
		*timeRemaining = resetTime
		return true
	}
	return false
}

func HasHitTime(timeRemaining *float32, deltaTime float32) bool {
	*timeRemaining = *timeRemaining - deltaTime
	return *timeRemaining <= 0
}
