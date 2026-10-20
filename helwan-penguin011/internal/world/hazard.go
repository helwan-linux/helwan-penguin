package world

// HazardKind identifies hazards with different gameplay behavior.
type HazardKind int

const (
	HazardPit HazardKind = iota
	HazardIceSpikes
	HazardThorns
	HazardLava
	HazardQuicksand
	HazardLaser
	HazardEnergy
	HazardFallingRock
	HazardWind
)

// Hazard may be static, timed, or move vertically (falling rocks).
type Hazard struct {
	X, Y, W, H float64
	Kind       HazardKind
	Period     float64
	ActiveFor  float64
	Phase      float64
	BaseY      float64
	Speed      float64
	Falling    bool
}

// ActiveAt returns whether a timed hazard is dangerous at the given frame.
func (h Hazard) ActiveAt(frame float64) bool {
	if h.Period <= 0 || h.ActiveFor <= 0 {
		return true
	}
	return mathMod(frame+h.Phase, h.Period) < h.ActiveFor
}

func mathMod(value, period float64) float64 {
	value = value - float64(int(value/period))*period
	if value < 0 {
		value += period
	}
	return value
}
