// Package main 👍
package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	collider "github.com/vcscsvcscs/ebiten-collider"
)

// vars
var (
	WindowWidth  = 640 * 2
	WindowHeight = 480 * 2

	player *Player
	wall   *collider.RectangleShape
	wall2  *collider.RectangleShape
	obs    *collider.CircleShape
	obs2   *collider.CircleShape
	hash   *collider.SpatialHash
	cursor *collider.PointShape

	ErrNormalExit = errors.New("normal exit")
)

// Player is the moveable shape
type Player struct {
	Bounds *collider.CircleShape
	// Bounds *collider.RectangleShape
	Speed float64
}

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ErrNormalExit
	}

	if inpututil.KeyPressDuration(ebiten.KeyLeft) > 0 || inpututil.KeyPressDuration(ebiten.KeyH) > 0 {
		player.Bounds.Move(-player.Speed, 0)
	}
	if inpututil.KeyPressDuration(ebiten.KeyRight) > 0 || inpututil.KeyPressDuration(ebiten.KeyN) > 0 {
		player.Bounds.Move(player.Speed, 0)
	}
	if inpututil.KeyPressDuration(ebiten.KeyUp) > 0 || inpututil.KeyPressDuration(ebiten.KeyC) > 0 {
		player.Bounds.Move(0, -player.Speed)
	}
	if inpututil.KeyPressDuration(ebiten.KeyDown) > 0 || inpututil.KeyPressDuration(ebiten.KeyT) > 0 {
		player.Bounds.Move(0, player.Speed)
	}

	cx, cy := ebiten.CursorPosition()
	cursor.MoveTo(float64(cx), float64(cy))

	collisions := hash.CheckCollisions(player.Bounds)
	for _, collision := range collisions {
		sep := collision.SeparatingVector
		player.Bounds.Move(sep.X, sep.Y)

		log.Println(collision.Other.GetParent())
		// player.Bounds.Move(sep.X/2, sep.Y/2)
		// collision.Other.Move(-sep.X/2, -sep.Y/2)
	}

	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	hash.Draw(screen)

	red := color.RGBA{255, 0, 0, 128}
	green := color.RGBA{0, 255, 0, 128}
	_ = green

	vector.DrawFilledCircle(
		screen,
		float32(player.Bounds.Pos.X),
		float32(player.Bounds.Pos.Y),
		float32(player.Bounds.Radius),
		red, true)
	// vector.DrawFilledRect(
	// 	screen,
	// 	player.Bounds.Pos.X-player.Bounds.Width/2,
	// 	player.Bounds.Pos.Y-player.Bounds.Height/2,
	// 	player.Bounds.Width,
	// 	player.Bounds.Height,
	// 	red)

	vector.DrawFilledCircle(
		screen,
		float32(obs.Pos.X),
		float32(obs.Pos.Y),
		float32(obs.Radius),
		red, true)
	vector.DrawFilledCircle(
		screen,
		float32(obs2.Pos.X),
		float32(obs2.Pos.Y),
		float32(obs2.Radius),
		red, true)

	vector.DrawFilledRect(
		screen,
		float32(wall.Pos.X-wall.Width/2),
		float32(wall.Pos.Y-wall.Height/2),
		float32(wall.Width),
		float32(wall.Height),
		red, true)
	vector.DrawFilledRect(
		screen,
		float32(wall2.Pos.X-wall2.Width/2),
		float32(wall2.Pos.Y-wall2.Height/2),
		float32(wall2.Width),
		float32(wall2.Height),
		red, true)

	vector.DrawFilledCircle(
		screen,
		float32(cursor.Pos.X),
		float32(cursor.Pos.Y),
		5,
		red, true)
}

// Layout sets window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if WindowWidth != outsideWidth || WindowHeight != outsideHeight {
		log.Println("resize", outsideWidth, outsideHeight)
		WindowWidth = outsideWidth
		WindowHeight = outsideHeight
	}
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Collisions example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	x := float64(WindowWidth)/2 - 64/2
	y := float64(WindowHeight)/2 - 64/2

	hash = collider.NewSpatialHash(128)
	player = &Player{
		Bounds: hash.NewCircleShape(x, y, 32),
		// Bounds: hash.NewRectangleShape(x, y, 64, 64),
		Speed: 10,
	}
	player.Bounds.SetParent("I'm the player")

	wall = hash.NewRectangleShape(
		x,
		y-16,
		128,
		128,
	)
	wall.SetParent("Wall 1")
	wall2 = hash.NewRectangleShape(
		x+128,
		y-16,
		128,
		128*2,
	)
	wall2.SetParent("Wall 2")

	obs = hash.NewCircleShape(
		x+96,
		y+256+64,
		32)
	obs.SetParent("Circle 1")
	obs2 = hash.NewCircleShape(
		x+128,
		y+256,
		64)
	obs2.SetParent("Circle 2")

	cursor = hash.NewPointShape(0, 0)
	cursor.SetParent("Cursor")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
