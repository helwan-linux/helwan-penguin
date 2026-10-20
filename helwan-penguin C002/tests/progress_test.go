package tests

import (
	"github.com/helwan-linux/helwan-penguin/internal/progress"
	"os"
	"path/filepath"
	"testing"
)

func TestProgressSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.json")
	want := progress.Progress{UnlockedLevel: 3, CompletedLevels: []int{1, 2}, TotalStars: 7}
	if err := progress.Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := progress.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.UnlockedLevel != want.UnlockedLevel || got.TotalStars != want.TotalStars || len(got.CompletedLevels) != 2 {
		t.Fatalf("unexpected progress: %+v", got)
	}
	_ = os.Remove(path)
}
