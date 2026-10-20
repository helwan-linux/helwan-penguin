
package progress

type Progress struct {
	UnlockedLevel   int   `json:"unlocked_level"`
	CompletedLevels []int `json:"completed_levels"`
	TotalStars      int   `json:"total_stars"`
}

func Default() Progress {
	return Progress{
		UnlockedLevel:   1,
		CompletedLevels: []int{},
		TotalStars:      0,
	}
}

func (p Progress) IsUnlocked(level int) bool {
	return level >= 1 && level <= 35 && level <= p.UnlockedLevel
}

func (p Progress) IsCompleted(level int) bool {
	for _, completed := range p.CompletedLevels {
		if completed == level {
			return true
		}
	}
	return false
}
