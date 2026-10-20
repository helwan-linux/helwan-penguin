package levels

type LevelInfo struct {
	Number      int
	World       int
	Name        string
	Description string
}

type PlatformSpec struct {
	X float64
	Y float64
	W float64
}

func All() []LevelInfo {
	worlds := []struct {
		name        string
		description string
	}{
		{"Snowy Beginnings", "Snow and frozen ground"},
		{"Frozen Peaks", "High icy mountains"},
		{"Wild Jungle", "Dangerous jungle paths"},
		{"Lost Desert", "Ancient ruins and hot dunes"},
		{"Sky Islands", "Floating islands above the clouds"},
		{"Shadow Castle", "Dark halls and deadly traps"},
		{"Crystal Dimension", "The final crystal world"},
	}

	result := make([]LevelInfo, 0, 35)

	for world := 1; world <= 7; world++ {
		for stage := 1; stage <= 5; stage++ {
			result = append(result, LevelInfo{
				Number:      (world-1)*5 + stage,
				World:       world,
				Name:        worlds[world-1].name,
				Description: worlds[world-1].description,
			})
		}
	}

	return result
}

func Route(worldNumber, stage int) []PlatformSpec {
	if worldNumber < 1 {
		worldNumber = 1
	}
	if worldNumber > 7 {
		worldNumber = 7
	}
	if stage < 1 {
		stage = 1
	}
	if stage > 5 {
		stage = 5
	}

	var layouts [][]PlatformSpec

	switch worldNumber {
	case 1:
		layouts = world01Layouts()
	case 2:
		layouts = world02Layouts()
	case 3:
		layouts = world03Layouts()
	case 4:
		layouts = world04Layouts()
	case 5:
		layouts = world05Layouts()
	case 6:
		layouts = world06Layouts()
	case 7:
		layouts = world07Layouts()
	}

	return append([]PlatformSpec(nil), layouts[stage-1]...)
}
