
package progress

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const MaxLevel = 35

func SavePath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, "helwan-penguin", "progress.json"), nil
}

func resolvePath(paths []string) (string, error) {
	if len(paths) > 0 && paths[0] != "" {
		return paths[0], nil
	}
	return SavePath()
}

func Load(paths ...string) (Progress, error) {
	path, err := resolvePath(paths)
	if err != nil {
		return Default(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Default(), err
	}

	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		return Default(), err
	}

	if p.UnlockedLevel < 1 {
		p.UnlockedLevel = 1
	}
	if p.UnlockedLevel > MaxLevel {
		p.UnlockedLevel = MaxLevel
	}
	if p.CompletedLevels == nil {
		p.CompletedLevels = []int{}
	}
	if p.TotalStars < 0 {
		p.TotalStars = 0
	}

	return p, nil
}

func Save(args ...interface{}) error {
	var path string
	var p Progress

	switch len(args) {
	case 1:
		var ok bool
		p, ok = args[0].(Progress)
		if !ok {
			return os.ErrInvalid
		}

		var err error
		path, err = SavePath()
		if err != nil {
			return err
		}

	case 2:
		var ok bool
		path, ok = args[0].(string)
		if !ok || path == "" {
			return os.ErrInvalid
		}

		p, ok = args[1].(Progress)
		if !ok {
			return os.ErrInvalid
		}

	default:
		return os.ErrInvalid
	}

	if p.UnlockedLevel < 1 {
		p.UnlockedLevel = 1
	}
	if p.UnlockedLevel > MaxLevel {
		p.UnlockedLevel = MaxLevel
	}
	if p.CompletedLevels == nil {
		p.CompletedLevels = []int{}
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	temp := path + ".tmp"
	if err := os.WriteFile(temp, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(temp, path); err != nil {
		_ = os.Remove(temp)
		return err
	}

	return nil
}

func MarkCompleted(p *Progress, level int, stars int) {
	if p == nil || level < 1 || level > MaxLevel {
		return
	}

	for _, completed := range p.CompletedLevels {
		if completed == level {
			return
		}
	}

	p.CompletedLevels = append(p.CompletedLevels, level)

	if p.UnlockedLevel <= level && level < MaxLevel {
		p.UnlockedLevel = level + 1
	}
	if level == MaxLevel && p.UnlockedLevel < MaxLevel {
		p.UnlockedLevel = MaxLevel
	}
	if stars > 0 {
		p.TotalStars += stars
	}
}
