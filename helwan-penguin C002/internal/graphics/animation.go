package graphics

import "math"

// PenguinAnimation controls the penguin's movement animation.
type PenguinAnimation struct {
	Time    float64
	Moving  bool
	Jumping bool
	Falling bool
	Facing  int
}

// NewPenguinAnimation creates a new animation controller.
func NewPenguinAnimation() PenguinAnimation {
	return PenguinAnimation{
		Facing: 1,
	}
}

// Update advances the animation state.
func (a *PenguinAnimation) Update(moving, jumping, falling bool) {
	a.Moving = moving
	a.Jumping = jumping
	a.Falling = falling

	if a.Facing == 0 {
		a.Facing = 1
	}

	a.Time += 1.0 / 60.0
}

// WalkOffset returns the vertical movement of the walking animation.
func (a *PenguinAnimation) WalkOffset() float32 {
	if !a.Moving || a.Jumping || a.Falling {
		return 0
	}

	return float32(math.Sin(a.Time*12) * 3)
}

// WingAngle returns the current wing animation angle.
func (a *PenguinAnimation) WingAngle() float64 {
	if a.Jumping {
		return -0.7
	}

	if a.Falling {
		return 0.7
	}

	if a.Moving {
		return math.Sin(a.Time*12) * 0.35
	}

	return math.Sin(a.Time*2) * 0.05
}

// FootOffset returns the movement offset for the feet.
func (a *PenguinAnimation) FootOffset() float32 {
	if !a.Moving || a.Jumping || a.Falling {
		return 0
	}

	return float32(math.Sin(a.Time*12) * 4)
}
