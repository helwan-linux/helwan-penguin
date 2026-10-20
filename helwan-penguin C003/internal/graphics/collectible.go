package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

func DrawCollectible(screen *ebiten.Image, c world.Collectible, cameraX float64) {
	if c.Taken {
		return
	}

	ebitenutil.DrawCircle(
		screen,
		c.X-cameraX,
		c.Y,
		9,
		color.RGBA{R: 255, G: 194, B: 40, A: 255},
	)
}
