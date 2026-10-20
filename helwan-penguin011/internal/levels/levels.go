package levels

type LevelInfo struct {
	World, Number int
	Name          string
	Difficulty    int
}

func All() []LevelInfo {
	worlds := []string{"Snowy Beginnings", "Frozen Peaks", "Wild Jungle", "Lost Desert", "Sky Islands", "Shadow Castle"}
	out := make([]LevelInfo, 0, 30)
	for wi, name := range worlds {
		for n := 1; n <= 5; n++ {
			number := wi*5 + n
			out = append(out, LevelInfo{World: wi + 1, Number: number, Name: name, Difficulty: number})
		}
	}
	return out
}
