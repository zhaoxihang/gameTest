package main

import (
	"gameTest/core"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	ScreenWidth  = 420 //初始画布宽
	ScreenHeight = 600 //初始画布高
)

func main() {
	g, err := core.NewGame(ScreenWidth, ScreenHeight)
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("GameDemo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
