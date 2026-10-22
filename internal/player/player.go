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

	// Difficulty settings.
	MoveSpeed    float64
	JumpStrength float64
	Gravity      float64
	MaxFallSpeed float64
	JumpAssist   float64
}

func New(x, y float64) Player {
	p := Player{
		X:            x,
		Y:            y,
		Width:        30,
		Height:       38,
		MoveSpeed:    3.4,
		JumpStrength: -10.8,
		Gravity:      0.48,
		MaxFallSpeed: 12,
		JumpAssist:   0,
	}
	return p
}

// SetDifficulty configures movement for easy, normal, or hard mode.
// Accepted values: "easy", "normal", "hard".
func (p *Player) SetDifficulty(difficulty string) {
	switch difficulty {
	case "easy":
		p.MoveSpeed = 3.6
		p.JumpStrength = -12.0
		p.Gravity = 0.40
		p.MaxFallSpeed = 10
		p.JumpAssist = 5

	case "hard":
		p.MoveSpeed = 3.2
		p.JumpStrength = -10.0
		p.Gravity = 0.55
		p.MaxFallSpeed = 13
		p.JumpAssist = 0

	default: // normal
		p.MoveSpeed = 3.4
		p.JumpStrength = -10.8
		p.Gravity = 0.48
		p.MaxFallSpeed = 12
		p.JumpAssist = 0
	}
}

func (p *Player) Update() {
	// Keep safe defaults for players created without SetDifficulty.
	if p.MoveSpeed == 0 {
		p.MoveSpeed = 3.4
	}
	if p.JumpStrength == 0 {
		p.JumpStrength = -10.8
	}
	if p.Gravity == 0 {
		p.Gravity = 0.48
	}
	if p.MaxFallSpeed == 0 {
		p.MaxFallSpeed = 12
	}

	p.VX = 0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) ||
		ebiten.IsKeyPressed(ebiten.KeyA) {
		p.VX = -p.MoveSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.IsKeyPressed(ebiten.KeyD) {
		p.VX = p.MoveSpeed
	}

	jumpPressed := ebiten.IsKeyPressed(ebiten.KeySpace) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) ||
		ebiten.IsKeyPressed(ebiten.KeyW)

	if jumpPressed && p.OnGround {
		p.VY = p.JumpStrength
		p.OnGround = false
	}

	p.VY += p.Gravity

	if p.VY > p.MaxFallSpeed {
		p.VY = p.MaxFallSpeed
	}

	p.X += p.VX
	p.Y += p.VY

	if p.X < 0 {
		p.X = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image, cameraX float64) {
	x, y := p.X-cameraX, p.Y

	ebitenutil.DrawRect(screen, x+4, y+7, 22, 28,
		color.RGBA{R: 35, G: 45, B: 59, A: 255})
	ebitenutil.DrawRect(screen, x+8, y+12, 14, 19,
		color.RGBA{R: 245, G: 248, B: 250, A: 255})
	ebitenutil.DrawRect(screen, x+7, y, 16, 13,
		color.RGBA{R: 35, G: 45, B: 59, A: 255})
	ebitenutil.DrawRect(screen, x+18, y+5, 7, 4,
		color.RGBA{R: 255, G: 174, B: 45, A: 255})
	ebitenutil.DrawRect(screen, x+3, y+33, 11, 5,
		color.RGBA{R: 255, G: 174, B: 45, A: 255})
	ebitenutil.DrawRect(screen, x+17, y+33, 11, 5,
		color.RGBA{R: 255, G: 174, B: 45, A: 255})
}
