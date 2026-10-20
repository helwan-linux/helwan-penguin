package tests

import (
	"github.com/helwan-linux/helwan-penguin/internal/player"
	"testing"
)

func TestNewPlayerHasExpectedBounds(t *testing.T) {
	p := player.New(10, 20)
	if p.X != 10 || p.Y != 20 || p.Width <= 0 || p.Height <= 0 {
		t.Fatalf("unexpected player: %+v", p)
	}
}
