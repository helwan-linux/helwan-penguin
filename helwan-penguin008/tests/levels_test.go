package tests

import (
	"github.com/helwan-linux/helwan-penguin/internal/levels"
	"testing"
)

func TestThirtyLevelsAreDefined(t *testing.T) {
	if got := len(levels.All()); got != 30 {
		t.Fatalf("want 30 levels, got %d", got)
	}
}
