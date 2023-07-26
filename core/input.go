package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input represents the current key states.
type Input struct {
	//鼠标状态
	mouseState mouseState
	//鼠标位置的x轴
	mouseInitPosX int
	//鼠标位置的y轴
	mouseInitPosY int
	//鼠标移动方向
	mouseDir Dir
}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}

// Update updates the current input states.
func (i *Input) Update() {
	//根据鼠标的状态改变动画
	switch i.mouseState {
	case mouseStateNone: //空状态
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			i.mouseInitPosX = x
			i.mouseInitPosY = y
			i.mouseState = mouseStatePressing
		}
	case mouseStatePressing:
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			dx := x - i.mouseInitPosX
			dy := y - i.mouseInitPosY
			d, ok := vecToDir(dx, dy)
			if !ok {
				i.mouseState = mouseStateNone
				break
			}
			i.mouseDir = d
			i.mouseState = mouseStateSettled
		}
	case mouseStateSettled:
		i.mouseState = mouseStateNone
	}
}

// Dir returns a currently pressed direction.
// Dir returns false if no direction key is pressed.
func (i *Input) Dir() (Dir, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return DirUp, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return DirLeft, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return DirRight, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return DirDown, true
	}
	if i.mouseState == mouseStateSettled {
		return i.mouseDir, true
	}
	return 0, false
}
