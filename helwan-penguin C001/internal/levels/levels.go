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
		"Crystal Dimension",
	}

	out := make([]LevelInfo, 0, 35)

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
	case 7:
		routes = world07Layouts()
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

func world07Layouts() [][]PlatformSpec {
	return [][]PlatformSpec{
		{
			{X: 0, Y: 440, W: 230},
			{X: 285, Y: 390, W: 110},
			{X: 450, Y: 345, W: 110},
			{X: 620, Y: 390, W: 120},
			{X: 800, Y: 440, W: 230},
			{X: 1100, Y: 390, W: 120},
		},
		{
			{X: 0, Y: 440, W: 190},
			{X: 250, Y: 385, W: 100},
			{X: 405, Y: 325, W: 100},
			{X: 565, Y: 375, W: 100},
			{X: 725, Y: 315, W: 110},
			{X: 900, Y: 440, W: 220},
			{X: 1200, Y: 390, W: 130},
		},
		{
			{X: 0, Y: 440, W: 180},
			{X: 240, Y: 390, W: 90},
			{X: 390, Y: 330, W: 90},
			{X: 540, Y: 275, W: 90},
			{X: 700, Y: 340, W: 100},
			{X: 870, Y: 390, W: 100},
			{X: 1050, Y: 440, W: 200},
			{X: 1330, Y: 385, W: 130},
		},
		{
			{X: 0, Y: 440, W: 170},
			{X: 230, Y: 385, W: 85},
			{X: 375, Y: 325, W: 85},
			{X: 520, Y: 270, W: 85},
			{X: 670, Y: 315, W: 85},
			{X: 820, Y: 365, W: 90},
			{X: 980, Y: 410, W: 100},
			{X: 1150, Y: 440, W: 180},
			{X: 1410, Y: 385, W: 140},
		},
		{
			{X: 0, Y: 440, W: 180},
			{X: 245, Y: 380, W: 85},
			{X: 390, Y: 315, W: 85},
			{X: 540, Y: 260, W: 85},
			{X: 690, Y: 310, W: 85},
			{X: 840, Y: 360, W: 85},
			{X: 990, Y: 310, W: 85},
			{X: 1140, Y: 365, W: 100},
			{X: 1320, Y: 440, W: 240},
			{X: 1640, Y: 385, W: 150},
		},
	}
}
