package graphics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/helwan-linux/helwan-penguin/internal/world"
)

func DrawHazard(screen *ebiten.Image, h world.Hazard, cameraX float64) {
	x := h.X - cameraX
	if x < -h.W-20 || x > 980 || h.W <= 0 || h.H <= 0 {
		return
	}

	switch h.Kind {
	case world.HazardIceSpikes:
		drawIceSpikes(screen, x, h.Y, h.W, h.H)
	case world.HazardThorns:
		drawThorns(screen, x, h.Y, h.W, h.H)
	case world.HazardLava:
		drawLava(screen, x, h.Y, h.W, h.H)
	case world.HazardQuicksand:
		drawQuicksand(screen, x, h.Y, h.W, h.H)
	case world.HazardLaser:
		drawLaser(screen, x, h.Y, h.W, h.H)
	case world.HazardEnergy:
		drawEnergy(screen, x, h.Y, h.W, h.H)
	case world.HazardFallingRock:
		drawFallingRock(screen, x, h.Y, h.W, h.H)
	case world.HazardWind:
		drawWind(screen, x, h.Y, h.W, h.H)
	default:
		drawRock(screen, x, h.Y, h.W, h.H)
	}
}

func hazardRect(s *ebiten.Image, x, y, w, h float64, c color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	vector.DrawFilledRect(s, float32(x), float32(y), float32(w), float32(h), c, true)
}

func hazardCircle(s *ebiten.Image, x, y, r float64, c color.Color) {
	if r <= 0 {
		return
	}
	vector.DrawFilledCircle(s, float32(x), float32(y), float32(r), c, true)
}

func hazardTriangle(s *ebiten.Image, x, y, w, h float64, c color.Color) {
	// Triangle made from horizontal slices, avoiding unsupported path APIs.
	if w <= 0 || h <= 0 {
		return
	}
	slices := int(math.Max(3, math.Min(16, h)))
	for i := 0; i < slices; i++ {
		t := float64(i+1) / float64(slices)
		sliceW := w * t
		sliceX := x + (w-sliceW)/2
		sliceY := y + h*float64(i)/float64(slices)
		hazardRect(s, sliceX, sliceY, sliceW, h/float64(slices)+1, c)
	}
}

func drawIceSpikes(s *ebiten.Image, x, y, w, h float64) {
	base := color.RGBA{30, 115, 165, 255}
	ice := color.RGBA{155, 235, 255, 255}
	shine := color.RGBA{235, 255, 255, 255}

	hazardRect(s, x, y+h*0.82, w, h*0.18, base)
	count := int(math.Max(1, math.Floor(w/15)))
	step := w / float64(count)

	for i := 0; i < count; i++ {
		spikeH := h * (0.60 + float64(i%2)*0.18)
		hazardTriangle(s, x+float64(i)*step, y+h*0.82-spikeH, step, spikeH, ice)
		hazardRect(s, x+float64(i)*step+step*0.43, y+h*0.82-spikeH*0.70, step*0.08, spikeH*0.35, shine)
	}
}

func drawThorns(s *ebiten.Image, x, y, w, h float64) {
	dark := color.RGBA{20, 65, 30, 255}
	green := color.RGBA{35, 115, 45, 255}
	light := color.RGBA{105, 210, 75, 255}

	hazardRect(s, x, y+h*0.72, w, h*0.28, dark)
	count := int(math.Max(1, math.Floor(w/14)))
	step := w / float64(count)

	for i := 0; i < count; i++ {
		spikeH := h * (0.50 + float64(i%2)*0.18)
		hazardTriangle(s, x+float64(i)*step, y+h*0.72-spikeH, step, spikeH, green)
		hazardRect(s, x+float64(i)*step+step*0.45, y+h*0.72-spikeH*0.60, step*0.07, spikeH*0.22, light)
	}
}

func drawLava(s *ebiten.Image, x, y, w, h float64) {
	rock := color.RGBA{65, 35, 40, 255}
	red := color.RGBA{210, 40, 15, 255}
	orange := color.RGBA{255, 100, 15, 255}
	yellow := color.RGBA{255, 205, 55, 255}

	hazardRect(s, x, y, w, h, rock)
	hazardRect(s, x, y+h*0.18, w, h*0.55, red)

	count := int(math.Max(2, math.Floor(w/18)))
	step := w / float64(count)
	for i := 0; i < count; i++ {
		cx := x + (float64(i)+0.5)*step
		wave := math.Sin(float64(i)*1.7) * h * 0.09
		hazardCircle(s, cx, y+h*0.30+wave, step*0.32, orange)
		hazardRect(s, cx-step*0.20, y+h*0.48+wave, step*0.40, h*0.12, yellow)
	}
	hazardRect(s, x, y, w, math.Max(2, h*0.06), orange)
}

func drawQuicksand(s *ebiten.Image, x, y, w, h float64) {
	sand := color.RGBA{175, 115, 45, 255}
	light := color.RGBA{240, 195, 95, 255}
	dark := color.RGBA{115, 70, 35, 255}

	hazardRect(s, x, y, w, h, dark)
	hazardRect(s, x, y+h*0.12, w, h*0.88, sand)

	count := int(math.Max(2, math.Floor(w/22)))
	step := w / float64(count)
	for i := 0; i < count; i++ {
		cx := x + (float64(i)+0.5)*step
		offset := math.Sin(float64(i)*1.8) * h * 0.06
		hazardCircle(s, cx, y+h*0.20+offset, math.Max(2, step*0.23), light)
		hazardRect(s, cx-step*0.18, y+h*0.45+offset, step*0.36, math.Max(1, h*0.05), dark)
	}
}

func drawLaser(s *ebiten.Image, x, y, w, h float64) {
	dark := color.RGBA{70, 10, 30, 255}
	red := color.RGBA{255, 25, 65, 255}
	white := color.RGBA{255, 225, 235, 255}

	hazardRect(s, x, y, w, h, dark)
	hazardRect(s, x, y+h*0.25, w, h*0.5, red)
	hazardRect(s, x, y+h*0.40, w, h*0.20, white)
	hazardCircle(s, x, y+h*0.5, h*0.42, red)
	hazardCircle(s, x+w, y+h*0.5, h*0.42, red)
	hazardCircle(s, x, y+h*0.5, h*0.19, white)
	hazardCircle(s, x+w, y+h*0.5, h*0.19, white)
}

func drawEnergy(s *ebiten.Image, x, y, w, h float64) {
	blue := color.RGBA{15, 80, 180, 255}
	cyan := color.RGBA{45, 210, 255, 255}
	white := color.RGBA{220, 255, 255, 255}

	hazardRect(s, x, y, w, h, blue)
	hazardRect(s, x+w*0.20, y, w*0.60, h, cyan)

	count := int(math.Max(2, math.Floor(w/14)))
	step := w / float64(count)
	for i := 0; i < count; i++ {
		cx := x + (float64(i)+0.5)*step
		hazardCircle(s, cx, y+h*0.5, math.Max(1, math.Min(step, h)*0.15), white)
	}
}

func drawFallingRock(s *ebiten.Image, x, y, w, h float64) {
	dark := color.RGBA{55, 48, 50, 255}
	rock := color.RGBA{115, 100, 90, 255}
	light := color.RGBA{165, 145, 125, 255}

	r := math.Min(w, h)
	hazardCircle(s, x+w*0.5, y+h*0.52, r*0.48, dark)
	hazardCircle(s, x+w*0.47, y+h*0.45, r*0.39, rock)
	hazardCircle(s, x+w*0.35, y+h*0.35, r*0.13, light)
	hazardCircle(s, x+w*0.65, y+h*0.60, r*0.10, dark)
}

func drawWind(s *ebiten.Image, x, y, w, h float64) {
	pale := color.RGBA{170, 225, 255, 170}
	white := color.RGBA{235, 250, 255, 230}
	blue := color.RGBA{100, 185, 235, 200}

	for i := 0; i < 3; i++ {
		yy := y + h*(0.20+float64(i)*0.30)
		length := w * (0.65 + float64(i%2)*0.25)
		hazardRect(s, x, yy, length, math.Max(2, h*0.08), blue)
		hazardCircle(s, x+length, yy+h*0.04, math.Max(2, h*0.09), white)
		hazardRect(s, x+w*0.15, yy+h*0.10, length*0.55, math.Max(1, h*0.04), pale)
	}
}

func drawRock(s *ebiten.Image, x, y, w, h float64) {
	dark := color.RGBA{55, 50, 55, 255}
	rock := color.RGBA{105, 95, 100, 255}
	light := color.RGBA{160, 145, 145, 255}

	hazardRect(s, x, y, w, h, dark)
	hazardCircle(s, x+w*0.25, y+h*0.25, math.Min(w, h)*0.20, light)
	hazardCircle(s, x+w*0.70, y+h*0.65, math.Min(w, h)*0.18, rock)
}
