package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/helwan-linux/helwan-penguin/internal/game"
)

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Helwan Penguin")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
