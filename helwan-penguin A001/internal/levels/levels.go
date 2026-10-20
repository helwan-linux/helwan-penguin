
package levels

type LevelInfo struct {
	World      int
	Number     int
	Name       string
	Difficulty int
}

type PlatformSpec struct {
	X float64
	Y float64
	W float64
}

func All() []LevelInfo {
	worlds := []string{
		"Snowy Beginnings",
		"Frozen Peaks",
		"Wild Jungle",
		"Lost Desert",
		"Sky Islands",
		"Shadow Castle",
	}

	out := make([]LevelInfo, 0, 30)

	for wi, name := range worlds {
		for n := 1; n <= 5; n++ {
			number := wi*5 + n
			out = append(out, LevelInfo{
				World:      wi + 1,
				Number:     number,
				Name:       name,
				Difficulty: number,
			})
		}
	}

	return out
}

func Route(worldNumber, stage int) []PlatformSpec {
	var routes [][]PlatformSpec

	switch worldNumber {
	case 1:
		routes = world01Layouts()
	case 2:
		routes = world02Layouts()
	case 3:
		routes = world03Layouts()
	case 4:
		routes = world04Layouts()
	case 5:
		routes = world05Layouts()
	case 6:
		routes = world06Layouts()
	default:
		routes = world01Layouts()
	}

	if len(routes) == 0 {
		return nil
	}

	if stage < 1 {
		stage = 1
	}
	if stage > len(routes) {
		stage = len(routes)
	}

	return append([]PlatformSpec(nil), routes[stage-1]...)
}
