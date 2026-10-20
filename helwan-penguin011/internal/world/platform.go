package world

type PlatformKind int

const (
	PlatformNormal PlatformKind = iota
	PlatformIce
	PlatformMoving
	PlatformCrumbling
	PlatformConveyor
)

type Platform struct {
	X, Y, W, H       float64
	BaseX, BaseY     float64
	Kind             PlatformKind
	MoveRange, Speed float64
	MoveVertical     bool
	ConveyorSpeed    float64
	CrumbleTimer     float64
	Broken           bool
}
