package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawOverlay(screen *ebiten.Image, title, subtitle, footer string) {
	ebitenutil.DrawRect(screen, 180, 175, 600, 190,
		color.RGBA{R: 20, G: 35, B: 50, A: 235})

	ebitenutil.DebugPrintAt(screen, title, 390, 220)
	ebitenutil.DebugPrintAt(screen, subtitle, 260, 260)

	if footer != "" {
		ebitenutil.DebugPrintAt(screen, footer, 240, 295)
	}
}

func DrawInfo(screen *ebiten.Image, page string) {
	ebitenutil.DrawRect(screen, 120, 100, 720, 340,
		color.RGBA{R: 20, G: 35, B: 50, A: 245})

	if page == "about" {
		ebitenutil.DebugPrintAt(screen, "ABOUT HELWAN PENGUIN", 350, 140)
		ebitenutil.DebugPrintAt(screen, "A platform game by Helwan Linux", 285, 195)
		ebitenutil.DebugPrintAt(screen, "Arch Power, Helwan Simplicity", 300, 230)
		ebitenutil.DebugPrintAt(screen, "One Human Family, One Open Future.", 270, 265)
	} else {
		ebitenutil.DebugPrintAt(screen, "HOW TO PLAY", 400, 130)
		ebitenutil.DebugPrintAt(screen, "Move: A/D or Left/Right", 300, 170)
		ebitenutil.DebugPrintAt(screen, "Jump: Space or Up Arrow", 300, 200)
		ebitenutil.DebugPrintAt(screen, "Jump on enemies from above to defeat them", 245, 230)
		ebitenutil.DebugPrintAt(screen, "Collect stars and reach the trophy", 275, 260)
		ebitenutil.DebugPrintAt(screen, "P: Pause | N: New Game while paused", 280, 290)
		ebitenutil.DebugPrintAt(screen, "O: About | H: Help", 350, 320)
	}

	ebitenutil.DebugPrintAt(screen, "Press ESC or ENTER to return", 330, 390)
}
