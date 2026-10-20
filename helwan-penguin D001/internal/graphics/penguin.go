package graphics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	penguinBlack  = color.RGBA{R: 25, G: 35, B: 48, A: 255}
	penguinWhite  = color.RGBA{R: 245, G: 248, B: 250, A: 255}
	penguinOrange = color.RGBA{R: 255, G: 151, B: 35, A: 255}
	penguinEye    = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	penguinPupil  = color.RGBA{R: 15, G: 20, B: 25, A: 255}
)

// DrawPenguin draws the animated penguin.
// x and y represent the top-left corner of the player's hitbox.
func DrawPenguin(screen *ebiten.Image, x, y float32, animation *PenguinAnimation) {
	if animation == nil {
		return
	}

	facing := float32(animation.Facing)
	if facing == 0 {
		facing = 1
	}

	bob := animation.WalkOffset()
	if animation.Jumping {
		bob -= 2
	}
	if animation.Falling {
		bob += 1
	}

	// Keep the visual body aligned with the player's hitbox.
	cx := x + 15
	top := y + bob

	// Feet
	footOffset := animation.FootOffset()

	vector.DrawFilledCircle(
		screen,
		cx-7+footOffset,
		top+35,
		5,
		penguinOrange,
		true,
	)
	vector.DrawFilledCircle(
		screen,
		cx+7-footOffset,
		top+35,
		5,
		penguinOrange,
		true,
	)

	// Body
	vector.DrawFilledCircle(
		screen,
		cx,
		top+22,
		15,
		penguinBlack,
		true,
	)
	vector.DrawFilledRect(
		screen,
		cx-14,
		top+10,
		28,
		18,
		penguinBlack,
		true,
	)

	// Belly
	vector.DrawFilledCircle(
		screen,
		cx,
		top+24,
		10,
		penguinWhite,
		true,
	)
	vector.DrawFilledRect(
		screen,
		cx-9,
		top+22,
		18,
		8,
		penguinWhite,
		true,
	)

	// Head
	headY := top + 8
	vector.DrawFilledCircle(
		screen,
		cx,
		headY,
		12,
		penguinBlack,
		true,
	)

	// Eye
	eyeX := cx + facing*4
	vector.DrawFilledCircle(
		screen,
		eyeX,
		headY-1,
		3,
		penguinEye,
		true,
	)
	vector.DrawFilledCircle(
		screen,
		eyeX+facing,
		headY-1,
		1.4,
		penguinPupil,
		true,
	)

	// Beak
	vector.DrawFilledCircle(
		screen,
		cx+facing*11,
		headY+3,
		3.5,
		penguinOrange,
		true,
	)

	// Wing
	wingAngle := animation.WingAngle()
	if facing < 0 {
		wingAngle = -wingAngle
	}

	wingY := top + 22 + float32(math.Sin(wingAngle)*3)
	vector.DrawFilledCircle(
		screen,
		cx-facing*13,
		wingY,
		5,
		penguinBlack,
		true,
	)
}
