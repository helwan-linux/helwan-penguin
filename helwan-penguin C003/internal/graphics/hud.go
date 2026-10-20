package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawHUD(screen *ebiten.Image, worldNumber, level, attempts, maxAttempts, stars, score int, saveError string) {
	ebitenutil.DebugPrintAt(screen,
		"HELWAN PENGUIN | WORLD "+itoa(worldNumber), 18, 16)

	ebitenutil.DebugPrintAt(screen,
		"Level: "+itoa(level)+"/35  Attempts: "+
			itoa(attempts)+"/"+itoa(maxAttempts)+
			"  Stars: "+itoa(stars)+
			"  Score: "+itoa(score), 18, 38)

	if saveError != "" {
		ebitenutil.DebugPrintAt(screen, saveError, 18, 82)
	}
}

func DrawControls(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen,
		"Move: A/D or arrows | Jump: Space/Up | Pause: P", 18, 60)
}
