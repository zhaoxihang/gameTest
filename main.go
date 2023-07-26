package main

import (
	"gameTest/core"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(core.ScreenWidth, core.ScreenHeight)
	ebiten.SetWindowTitle("GameDemo")
	g, err := core.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
