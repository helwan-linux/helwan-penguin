package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

func DrawGoal(screen *ebiten.Image, w world.World, cameraX float64, finalLevel int) {
	if len(w.Platforms) == 0 {
		return
	}

	last := w.Platforms[len(w.Platforms)-1]
	x := w.EndX - cameraX

	width := 34.0
	top := last.Y - 58

	if w.Level == finalLevel {
		width = 48
		top = last.Y - 72
	}

	gold := color.RGBA{R: 255, G: 190, B: 35, A: 255}
	light := color.RGBA{R: 255, G: 225, B: 95, A: 255}

	ebitenutil.DrawRect(screen, x+width*0.28, top+24, width*0.44, 18, gold)
	ebitenutil.DrawRect(screen, x+width*0.40, top+42, width*0.20, 12, light)
	ebitenutil.DrawRect(screen, x+width*0.24, top+54, width*0.52, 6, gold)
	ebitenutil.DrawCircle(screen, x+width/2, top+18, width*0.40, gold)

	if w.Level == finalLevel {
		ebitenutil.DebugPrintAt(screen, "FINAL TROPHY", int(x-12), int(last.Y-92))
	}
}
