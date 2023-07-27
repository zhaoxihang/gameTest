package main

import (
	"gameTest/core"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	g, err := core.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("GameDemo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
