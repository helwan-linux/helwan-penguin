package tests

import (
	"testing"

	"github.com/helwan-linux/helwan-penguin/internal/levels"
)

func TestThirtyFiveLevelsAreDefined(t *testing.T) {
	got := levels.All()
	if len(got) != 35 {
		t.Fatalf("want 35 levels, got %d", len(got))
	}
}
