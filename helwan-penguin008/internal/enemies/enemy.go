package enemies

type Enemy struct {
	X, Y, Width, Height float64
	Active              bool
	VX                  float64
	MinX, MaxX          float64
	Boss                bool
	Health              int
}
