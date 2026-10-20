package graphics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawGhost(screen *ebiten.Image, x, y, width, height float64, body, head color.Color) {
	if width <= 0 || height <= 0 {
		return
	}

	fx, fy := float32(x), float32(y)
	w, h := float32(width), float32(height)

	// Shadow.
	vector.DrawFilledCircle(screen, fx+w*0.5, fy+h*0.94, w*0.38,
		color.RGBA{0, 0, 0, 55}, true)

	// Main ghost body.
	vector.DrawFilledCircle(screen, fx+w*0.5, fy+h*0.34, w*0.48, body, true)
	vector.DrawFilledRect(screen, fx+w*0.02, fy+h*0.34, w*0.96, h*0.48, body, true)

	// Wavy bottom edge.
	vector.DrawFilledCircle(screen, fx+w*0.20, fy+h*0.80, w*0.18, body, true)
	vector.DrawFilledCircle(screen, fx+w*0.50, fy+h*0.80, w*0.18, body, true)
	vector.DrawFilledCircle(screen, fx+w*0.80, fy+h*0.80, w*0.18, body, true)

	// Head highlight.
	vector.DrawFilledCircle(screen, fx+w*0.50, fy+h*0.20, float32(math.Max(1, width*0.12)), head, true)

	// Eyes.
	eyeY := fy + h*0.43
	eyeR := float32(math.Max(1.5, width*0.075))

	vector.DrawFilledCircle(screen, fx+w*0.35, eyeY, eyeR*1.7, color.RGBA{25, 30, 55, 255}, true)
	vector.DrawFilledCircle(screen, fx+w*0.65, eyeY, eyeR*1.7, color.RGBA{25, 30, 55, 255}, true)

	vector.DrawFilledCircle(screen, fx+w*0.35+eyeR*0.3, eyeY-eyeR*0.35, eyeR*0.55, color.White, true)
	vector.DrawFilledCircle(screen, fx+w*0.65+eyeR*0.3, eyeY-eyeR*0.35, eyeR*0.55, color.White, true)

	// Mouth.
	vector.DrawFilledCircle(screen, fx+w*0.50, fy+h*0.58, float32(math.Max(1, width*0.035)), color.RGBA{45, 20, 40, 255}, true)

	// Cheeks.
	vector.DrawFilledCircle(screen, fx+w*0.21, fy+h*0.55, float32(math.Max(1, width*0.045)), color.RGBA{255, 130, 170, 170}, true)
	vector.DrawFilledCircle(screen, fx+w*0.79, fy+h*0.55, float32(math.Max(1, width*0.045)), color.RGBA{255, 130, 170, 170}, true)
}
