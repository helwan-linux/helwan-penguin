package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

func DrawEnemy(screen *ebiten.Image, e world.Enemy, cameraX float64) {
	if !e.Active {
		return
	}

	x := e.X - cameraX
	if x < -e.Width || x > 960 {
		return
	}

	// Invisible ghosts cannot be seen.
	if e.Ghost && e.Invisible {
		return
	}

	// Draw the green stone left by a defeated world-nine ghost.
	if e.Stone {
		ebitenutil.DrawRect(
			screen,
			x,
			e.Y,
			e.Width,
			e.Height,
			color.RGBA{R: 25, G: 145, B: 55, A: 255},
		)

		if e.Width > 8 && e.Height > 8 {
			ebitenutil.DrawRect(
				screen,
				x+4,
				e.Y+4,
				e.Width-8,
				e.Height-8,
				color.RGBA{R: 100, G: 220, B: 100, A: 255},
			)
		}

		return
	}

	body := color.RGBA{R: 190, G: 55, B: 65, A: 255}
	head := color.RGBA{R: 220, G: 75, B: 80, A: 255}

	switch e.Kind {
	case world.EnemyRunner:
		body = color.RGBA{R: 215, G: 75, B: 45, A: 255}
		head = color.RGBA{R: 250, G: 115, B: 55, A: 255}

	case world.EnemyFlyer:
		body = color.RGBA{R: 45, G: 125, B: 190, A: 255}
		head = color.RGBA{R: 75, G: 185, B: 235, A: 255}

	case world.EnemyArmored:
		body = color.RGBA{R: 85, G: 95, B: 105, A: 255}
		head = color.RGBA{R: 155, G: 165, B: 175, A: 255}

	case world.EnemyGhost:
		body = color.RGBA{R: 115, G: 55, B: 175, A: 255}
		head = color.RGBA{R: 190, G: 135, B: 255, A: 255}
	}

	if e.Boss {
		body = color.RGBA{R: 90, G: 35, B: 145, A: 255}
		head = color.RGBA{R: 180, G: 105, B: 255, A: 255}
	}

	if e.HitFlash > 0 {
		body = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		head = color.RGBA{R: 255, G: 240, B: 120, A: 255}
	}

	DrawGhost(screen, x, e.Y, e.Width, e.Height, body, head)

	if e.Boss {
		ebitenutil.DebugPrintAt(
			screen,
			"Boss HP: "+itoa(e.Health),
			int(x),
			int(e.Y)-18,
		)
	} else if e.Ghost {
		ebitenutil.DebugPrintAt(
			screen,
			"Ghost HP: "+itoa(e.Health),
			int(x),
			int(e.Y)-18,
		)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	if n < 0 {
		return "-" + itoa(-n)
	}

	var b [20]byte
	i := len(b)

	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	return string(b[i:])
}
