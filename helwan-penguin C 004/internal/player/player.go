package player

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	X, Y          float64
	VX, VY        float64
	Width, Height float64
	OnGround      bool
}

func New(x, y float64) Player {
	return Player{X: x, Y: y, Width: 30, Height: 38}
}

func (p *Player) Update() {
	const speed, gravity, jump = 3.4, 0.48, -10.8

	p.VX = 0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.VX = -speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.VX = speed
	}

	if (ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.IsKeyPressed(ebiten.KeyW)) && p.OnGround {
		p.VY = jump
		p.OnGround = false
	}

	p.VY += gravity

	if p.VY > 12 {
		p.VY = 12
	}

	p.X += p.VX
	p.Y += p.VY

	if p.X < 0 {
		p.X = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image, cameraX float64) {
	x, y := p.X-cameraX, p.Y

	// A simple penguin assembled from shapes; no external art assets are needed yet.
	ebitenutil.DrawRect(screen, x+4, y+7, 22, 28, color.RGBA{R: 35, G: 45, B: 59, A: 255})
	ebitenutil.DrawRect(screen, x+8, y+12, 14, 19, color.RGBA{R: 245, G: 248, B: 250, A: 255})
	ebitenutil.DrawRect(screen, x+7, y, 16, 13, color.RGBA{R: 35, G: 45, B: 59, A: 255})
	ebitenutil.DrawRect(screen, x+18, y+5, 7, 4, color.RGBA{R: 255, G: 174, B: 45, A: 255})
	ebitenutil.DrawRect(screen, x+3, y+33, 11, 5, color.RGBA{R: 255, G: 174, B: 45, A: 255})
	ebitenutil.DrawRect(screen, x+17, y+33, 11, 5, color.RGBA{R: 255, G: 174, B: 45, A: 255})
}
