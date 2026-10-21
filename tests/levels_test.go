package tests

import (
	"testing"

	"github.com/helwan-linux/helwan-penguin/internal/levels"
)

func TestThirtyFiveLevelsAreDefined(t *testing.T) {
	got := levels.All()
	if len(got) != 45 {
		t.Fatalf("want 45 levels, got %d", len(got))
	}
}
