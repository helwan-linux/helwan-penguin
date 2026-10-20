package main

import (
	"bytes"
	_ "embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/helwan-linux/helwan-penguin/internal/game"
)

//go:embed icon.png
var windowIcon []byte

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Helwan Penguin")

	icon, _, err := image.Decode(bytes.NewReader(windowIcon))
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowIcon([]image.Image{icon})
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
