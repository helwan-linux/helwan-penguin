package graphics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	backgroundWidth  = 960
	backgroundHeight = 540
)

type Snowflake struct {
	X     float32
	Y     float32
	Size  float32
	Speed float32
	Drift float32
}

type SnowBackground struct {
	Frame      float64
	Snowflakes []Snowflake
}

type WorldBackground struct {
	Frame float64
	Snow  *SnowBackground
}

func NewWorldBackground() *WorldBackground {
	return &WorldBackground{
		Snow: NewSnowBackground(),
	}
}

func (b *WorldBackground) Update() {
	b.Frame++

	if b.Snow != nil {
		b.Snow.Update()
	}
}

func (b *WorldBackground) Draw(screen *ebiten.Image, cameraX float64, worldNumber int) {
	if worldNumber <= 1 {
		if b.Snow != nil {
			b.Snow.Draw(screen, cameraX)
			return
		}
	}

	switch worldNumber {
	case 2:
		drawGradientColors(screen,
			color.RGBA{R: 105, G: 190, B: 235, A: 255},
			color.RGBA{R: 190, G: 225, B: 150, A: 255})
		drawForest(screen, cameraX)
	case 3:
		drawGradientColors(screen,
			color.RGBA{R: 55, G: 24, B: 55, A: 255},
			color.RGBA{R: 225, G: 90, B: 35, A: 255})
		drawVolcano(screen, cameraX)
	case 4:
		drawGradientColors(screen,
			color.RGBA{R: 50, G: 145, B: 220, A: 255},
			color.RGBA{R: 250, G: 195, B: 105, A: 255})
		drawDesert(screen, cameraX)
	case 5:
		drawGradientColors(screen,
			color.RGBA{R: 35, G: 48, B: 72, A: 255},
			color.RGBA{R: 115, G: 130, B: 145, A: 255})
		drawIndustrial(screen, cameraX)
	default:
		drawGradientColors(screen,
			color.RGBA{R: 5, G: 8, B: 35, A: 255},
			color.RGBA{R: 45, G: 35, B: 100, A: 255})
		drawSpace(screen, cameraX)
	}
}

func NewSnowBackground() *SnowBackground {
	snow := &SnowBackground{
		Snowflakes: make([]Snowflake, 100),
	}

	for i := range snow.Snowflakes {
		snow.Snowflakes[i] = Snowflake{
			X:     float32((i*97 + 31) % backgroundWidth),
			Y:     float32((i*53 + 17) % backgroundHeight),
			Size:  1 + float32(i%4)*0.6,
			Speed: 0.7 + float32(i%5)*0.25,
			Drift: 0.3 + float32(i%4)*0.2,
		}
	}

	return snow
}

func (b *SnowBackground) Update() {
	b.Frame++

	for i := range b.Snowflakes {
		s := &b.Snowflakes[i]

		s.Y += s.Speed
		s.X += float32(math.Sin(b.Frame*0.025+float64(i))) * s.Drift

		if s.Y > backgroundHeight+5 {
			s.Y = -5
			s.X = float32((i*97 + 31) % backgroundWidth)
		}

		if s.X < -5 {
			s.X = backgroundWidth + 5
		} else if s.X > backgroundWidth+5 {
			s.X = -5
		}
	}
}

func (b *SnowBackground) Draw(screen *ebiten.Image, cameraX float64) {
	drawGradient(screen)

	for i := 0; i < 75; i++ {
		x := float32((i*137 + 31) % backgroundWidth)
		y := float32((i*53 + 17) % 245)

		vector.DrawFilledCircle(
			screen, x, y, 1.1,
			color.RGBA{R: 235, G: 249, B: 255, A: uint8(140 + i%4*25)},
			false,
		)
	}

	for i := 0; i < 5; i++ {
		x := float32(90+i*190) +
			float32(math.Sin(b.Frame*0.01+float64(i))*20)
		y := float32(105) +
			float32(math.Sin(b.Frame*0.015+float64(i))*12)

		vector.DrawFilledCircle(
			screen, x, y, 65,
			color.RGBA{R: 50, G: 220, B: 190, A: 18}, false,
		)
		vector.DrawFilledCircle(
			screen, x, y, 35,
			color.RGBA{R: 90, G: 245, B: 210, A: 16}, false,
		)
	}

	drawMountainRange(screen, cameraX*0.10, 335, 155, 170,
		color.RGBA{R: 72, G: 126, B: 168, A: 255},
		color.RGBA{R: 208, G: 235, B: 248, A: 255})

	drawMountainRange(screen, cameraX*0.22, 390, 135, 130,
		color.RGBA{R: 45, G: 97, B: 143, A: 255},
		color.RGBA{R: 224, G: 243, B: 252, A: 255})

	drawMountainRange(screen, cameraX*0.38, 435, 95, 90,
		color.RGBA{R: 29, G: 75, B: 112, A: 255},
		color.RGBA{R: 190, G: 225, B: 243, A: 255})

	for _, s := range b.Snowflakes {
		vector.DrawFilledCircle(
			screen, s.X, s.Y, s.Size,
			color.RGBA{R: 245, G: 251, B: 255, A: 225},
			false,
		)
	}
}

func drawGradient(screen *ebiten.Image) {
	drawGradientColors(screen,
		color.RGBA{R: 12, G: 30, B: 70, A: 255},
		color.RGBA{R: 105, G: 180, B: 215, A: 255})
}

func drawGradientColors(screen *ebiten.Image, top, bottom color.RGBA) {
	for y := 0; y < backgroundHeight; y++ {
		t := float64(y) / float64(backgroundHeight-1)

		c := color.RGBA{
			R: uint8(float64(top.R)*(1-t) + float64(bottom.R)*t),
			G: uint8(float64(top.G)*(1-t) + float64(bottom.G)*t),
			B: uint8(float64(top.B)*(1-t) + float64(bottom.B)*t),
			A: 255,
		}

		vector.DrawFilledRect(screen, 0, float32(y),
			backgroundWidth, 1, c, false)
	}
}

func drawMountainRange(
	screen *ebiten.Image,
	cameraX float64,
	baseY, peakHeight, spacing int,
	mountainColor, snowColor color.RGBA,
) {
	offset := int(math.Mod(cameraX, float64(spacing)))
	if offset < 0 {
		offset += spacing
	}

	for x := -spacing - offset; x < backgroundWidth+spacing; x += spacing {
		index := int(math.Abs(float64(x+offset)) / float64(spacing))
		variation := (index*37 + 19) % (peakHeight / 3)

		peakX := x + spacing/2
		peakY := baseY - peakHeight + variation

		var mountain vector.Path
		mountain.MoveTo(float32(x), float32(baseY))
		mountain.LineTo(float32(peakX), float32(peakY))
		mountain.LineTo(float32(x+spacing), float32(baseY))
		mountain.Close()

		vertices, indices := mountain.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vertices {
			vertices[i].ColorR = float32(mountainColor.R) / 255
			vertices[i].ColorG = float32(mountainColor.G) / 255
			vertices[i].ColorB = float32(mountainColor.B) / 255
			vertices[i].ColorA = float32(mountainColor.A) / 255
		}

		screen.DrawTriangles(vertices, indices, ebiten.NewImage(3, 3),
			&ebiten.DrawTrianglesOptions{})

		var cap vector.Path
		cap.MoveTo(float32(peakX-spacing/5), float32(peakY+peakHeight/4))
		cap.LineTo(float32(peakX), float32(peakY))
		cap.LineTo(float32(peakX+spacing/5), float32(peakY+peakHeight/4))
		cap.LineTo(float32(peakX), float32(peakY+peakHeight/6))
		cap.Close()

		capVertices, capIndices := cap.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range capVertices {
			capVertices[i].ColorR = float32(snowColor.R) / 255
			capVertices[i].ColorG = float32(snowColor.G) / 255
			capVertices[i].ColorB = float32(snowColor.B) / 255
			capVertices[i].ColorA = float32(snowColor.A) / 255
		}

		screen.DrawTriangles(capVertices, capIndices, ebiten.NewImage(3, 3),
			&ebiten.DrawTrianglesOptions{})
	}
}

func drawForest(screen *ebiten.Image, cameraX float64) {
	for i := 0; i < 15; i++ {
		x := float32(math.Mod(float64(i*117)-cameraX*0.25, 1200))
		if x < -30 {
			x += 1200
		}
		y := float32(300 + (i%3)*28)

		vector.DrawFilledRect(screen, x+12, y, 8, 180,
			color.RGBA{R: 95, G: 65, B: 38, A: 255}, false)
		vector.DrawFilledCircle(screen, x+16, y, 37,
			color.RGBA{R: 25, G: 105, B: 55, A: 255}, false)
		vector.DrawFilledCircle(screen, x-3, y+18, 27,
			color.RGBA{R: 35, G: 130, B: 65, A: 255}, false)
		vector.DrawFilledCircle(screen, x+34, y+18, 27,
			color.RGBA{R: 45, G: 145, B: 70, A: 255}, false)
	}

	vector.DrawFilledRect(screen, 0, 485, 960, 55,
		color.RGBA{R: 45, G: 105, B: 48, A: 255}, false)
}

func drawVolcano(screen *ebiten.Image, cameraX float64) {
	for i := 0; i < 7; i++ {
		x := float32(math.Mod(float64(i*190)-cameraX*0.15, 1300))
		if x < -100 {
			x += 1300
		}

		var mountain vector.Path
		mountain.MoveTo(x-90, 470)
		mountain.LineTo(x+20, float32(250+(i%3)*25))
		mountain.LineTo(x+130, 470)
		mountain.Close()

		vertices, indices := mountain.AppendVerticesAndIndicesForFilling(nil, nil)
		for j := range vertices {
			vertices[j].ColorR = 0.20
			vertices[j].ColorG = 0.12
			vertices[j].ColorB = 0.16
			vertices[j].ColorA = 1
		}
		screen.DrawTriangles(vertices, indices, ebiten.NewImage(3, 3),
			&ebiten.DrawTrianglesOptions{})

		vector.DrawFilledCircle(screen, x+20, float32(270+(i%3)*25), 15,
			color.RGBA{R: 255, G: 95, B: 20, A: 255}, false)
	}

	vector.DrawFilledRect(screen, 0, 480, 960, 60,
		color.RGBA{R: 100, G: 35, B: 30, A: 255}, false)
}

func drawDesert(screen *ebiten.Image, cameraX float64) {
	vector.DrawFilledCircle(screen, 790, 105, 43,
		color.RGBA{R: 255, G: 235, B: 145, A: 255}, false)

	for i := 0; i < 8; i++ {
		x := float32(math.Mod(float64(i*155)-cameraX*0.18, 1250))
		if x < -120 {
			x += 1250
		}

		vector.DrawFilledCircle(screen, x, 430, 110,
			color.RGBA{R: 215, G: 160, B: 85, A: 255}, false)
		vector.DrawFilledCircle(screen, x+90, 460, 100,
			color.RGBA{R: 235, G: 180, B: 100, A: 255}, false)
	}

	for i := 0; i < 5; i++ {
		x := float32(math.Mod(float64(i*245)-cameraX*0.30, 1200))
		if x < -30 {
			x += 1200
		}
		y := float32(345 + (i%2)*35)

		vector.DrawFilledRect(screen, x+12, y, 12, 115,
			color.RGBA{R: 35, G: 115, B: 65, A: 255}, false)
		vector.DrawFilledRect(screen, x-5, y+22, 18, 12,
			color.RGBA{R: 35, G: 115, B: 65, A: 255}, false)
		vector.DrawFilledRect(screen, x-5, y+22, 12, 42,
			color.RGBA{R: 35, G: 115, B: 65, A: 255}, false)
		vector.DrawFilledRect(screen, x+22, y+40, 20, 12,
			color.RGBA{R: 35, G: 115, B: 65, A: 255}, false)
		vector.DrawFilledRect(screen, x+30, y+40, 12, 32,
			color.RGBA{R: 35, G: 115, B: 65, A: 255}, false)
	}
}

func drawIndustrial(screen *ebiten.Image, cameraX float64) {
	for i := 0; i < 14; i++ {
		x := float32(math.Mod(float64(i*95)-cameraX*0.20, 1250))
		if x < -80 {
			x += 1250
		}

		height := float32(100 + (i%4)*35)
		vector.DrawFilledRect(screen, x, 455-height, 72, height,
			color.RGBA{R: 55, G: 65, B: 83, A: 255}, false)

		for row := 0; row < 3; row++ {
			for col := 0; col < 2; col++ {
				vector.DrawFilledRect(screen,
					x+10+float32(col*28),
					465-height+float32(row*28), 12, 13,
					color.RGBA{R: 240, G: 190, B: 75, A: 220}, false)
			}
		}
	}

	vector.DrawFilledRect(screen, 0, 475, 960, 65,
		color.RGBA{R: 50, G: 60, B: 72, A: 255}, false)
}

func drawSpace(screen *ebiten.Image, cameraX float64) {
	for i := 0; i < 115; i++ {
		x := float32(math.Mod(float64(i*137)-cameraX*0.08, 1100))
		if x < 0 {
			x += 1100
		}
		y := float32((i*53 + 17) % 430)
		size := float32(1 + i%3)

		vector.DrawFilledCircle(screen, x, y, size,
			color.RGBA{R: 220, G: 225, B: 255, A: 220}, false)
	}

	vector.DrawFilledCircle(screen, 740, 160, 65,
		color.RGBA{R: 125, G: 90, B: 205, A: 255}, false)
	vector.DrawFilledCircle(screen, 720, 140, 20,
		color.RGBA{R: 165, G: 145, B: 235, A: 255}, false)
	vector.DrawFilledCircle(screen, 200, 270, 32,
		color.RGBA{R: 60, G: 150, B: 215, A: 255}, false)
}
