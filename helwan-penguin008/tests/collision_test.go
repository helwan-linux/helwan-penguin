package tests

import (
	"github.com/helwan-linux/helwan-penguin/internal/player"
	"github.com/helwan-linux/helwan-penguin/internal/world"
	"testing"
)

func TestFirstLevelHasPlatformsAndCollectibles(t *testing.T) {
	w := world.FirstLevel()
	if len(w.Platforms) < 2 || len(w.Collectibles) == 0 {
		t.Fatalf("level is incomplete: %+v", w)
	}
	p := player.New(30, 400)
	p.VY = 1
	w.ResolvePlayer(&p)
	if !p.OnGround {
		t.Fatal("player should land on the starting platform")
	}
}
