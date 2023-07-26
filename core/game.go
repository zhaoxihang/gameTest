package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 420 //初始画布宽
	ScreenHeight = 600 //初始画布高
	boardSize    = 4
)

type Game struct {
	input      *Input
	board      *Board
	boardImage *ebiten.Image
}

func NewGame() (*Game, error) {
	g := &Game{
		input: NewInput(),
	}
	var err error
	g.board, err = NewBoard(boardSize)
	if err != nil {
		return g, err
	}
	return g, nil
}

// Update
// 更新更新游戏的逻辑状态
// 每一tick都会调用这个函数,Tick是逻辑更新的时间单位。默认值为1/60[S]，则默认每秒调用60次更新(即一个Ebiten游戏每秒调用60次)。
// UPDATE返回错误值。在这段代码中，更新总是返回nil。
// 通常，当更新函数返回非零错误时，Ebiten游戏暂停。
// 由于该程序从不返回非零错误，因此除非用户关闭窗口，否则Ebiten游戏永远不会停止。
func (g *Game) Update() error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	return nil
}

// Draw
// 每一帧都会调用这个函数
// 帧是渲染的时间单位，这取决于显示器的刷新率。如果监视器的刷新率为60[赫兹]，则每秒调用DRAW 60次。
// Draw接受一个参数屏幕，它是指向ebiten.Image的指针。
// 在ebiten中，所有图像(如从图像文件创建的图像、屏幕外图像(临时渲染目标)和屏幕)都表示为ebiten.Image对象。
// 屏幕参数是渲染的最终目的地。该窗口每帧显示屏幕的最终状态。
func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(g.board.Size())
	}
	//设置背景颜色
	screen.Fill(backgroundColor)
	//渲染棋盘
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	//游戏界面的大小
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	//棋盘的大小
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	//棋盘居中
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
