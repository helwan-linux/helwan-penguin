package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

func DrawPlatform(screen *ebiten.Image, p world.Platform, cameraX float64, worldNumber int) {
	if p.Broken {
		return
	}

	x := p.X - cameraX
	if x < -p.W || x > 960 {
		return
	}

	ground, top := platformColors(worldNumber)

	switch p.Kind {
	case world.PlatformIce:
		ground = color.RGBA{R: 120, G: 205, B: 238, A: 255}
		top = color.RGBA{R: 235, G: 250, B: 255, A: 255}
	case world.PlatformMoving:
		ground = color.RGBA{R: 130, G: 160, B: 180, A: 255}
		top = color.RGBA{R: 255, G: 220, B: 100, A: 255}
	case world.PlatformCrumbling:
		ground = color.RGBA{R: 155, G: 105, B: 65, A: 255}
		top = color.RGBA{R: 235, G: 165, B: 85, A: 255}
	case world.PlatformConveyor:
		ground = color.RGBA{R: 95, G: 105, B: 115, A: 255}
		top = color.RGBA{R: 245, G: 185, B: 55, A: 255}
	}

	ebitenutil.DrawRect(screen, x, p.Y, p.W, p.H, ground)
	ebitenutil.DrawRect(screen, x, p.Y, p.W, 6, top)
}

func platformColors(worldNumber int) (color.RGBA, color.RGBA) {
	switch worldNumber {
	case 2:
		return color.RGBA{R: 58, G: 105, B: 62, A: 255},
			color.RGBA{R: 145, G: 205, B: 105, A: 255}
	case 3:
		return color.RGBA{R: 105, G: 49, B: 45, A: 255},
			color.RGBA{R: 255, G: 125, B: 55, A: 255}
	case 4:
		return color.RGBA{R: 194, G: 151, B: 83, A: 255},
			color.RGBA{R: 255, G: 220, B: 135, A: 255}
	case 5:
		return color.RGBA{R: 65, G: 75, B: 91, A: 255},
			color.RGBA{R: 105, G: 205, B: 220, A: 255}
	case 6:
		return color.RGBA{R: 65, G: 48, B: 105, A: 255},
			color.RGBA{R: 165, G: 125, B: 255, A: 255}
	case 7:
		return color.RGBA{R: 70, G: 35, B: 55, A: 255},
			color.RGBA{R: 255, G: 90, B: 130, A: 255}
	default:
		return color.RGBA{R: 116, G: 173, B: 194, A: 255},
			color.RGBA{R: 245, G: 252, B: 255, A: 255}
	}
}
