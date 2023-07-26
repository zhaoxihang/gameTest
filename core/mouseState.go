package core

type mouseState int //鼠标状态

const (
	mouseStateNone     mouseState = iota //空状态
	mouseStatePressing                   //按下鼠标
	mouseStateSettled                    //鼠标复原
)
