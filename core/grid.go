package core

import (
	"gameTest/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"strconv"
)

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

const (
	maxMovingCount  = 5
	maxPoppingCount = 6
)
const (
	tileSize   = 80
	tileMargin = 4
)

var (
	gridImage = ebiten.NewImage(tileSize, tileSize)
)

// 初始化格子的动画
func init() {
	gridImage.Fill(color.White)
}

// 初始化字体
func init() {
	tt, err := opentype.Parse(fonts.Shangshou)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Grid struct {
	current GridData //当前格子

	next GridData //位移之后的格子 不会移动时为空

	movingCount       int //移动的步数
	startPoppingCount int //开始弹出的数字
	poppingCount      int //弹出计数
}

type GridData struct {
	value int //格子的数字
	x     int //x轴
	y     int //y轴
}

// NewGrid 初始化格子对象
func NewGrid(value int, x, y int) *Grid {
	return &Grid{
		current: GridData{
			value: value,
			x:     x,
			y:     y,
		},
		startPoppingCount: maxPoppingCount,
	}
}

// Pos 格子的位置
func (t *Grid) Pos() (int, int) {
	return t.current.x, t.current.y
}

// NextPos 下一步的位置
func (t *Grid) NextPos() (int, int) {
	return t.next.x, t.next.y
}

// Value 格子的数字
func (t *Grid) Value() int {
	return t.current.value
}

// NextValue 下一步格子的数字
func (t *Grid) NextValue() int {
	return t.next.value
}

// IsMoving 正在移动 移动的步数是否大于0
func (t *Grid) IsMoving() bool {
	return 0 < t.movingCount
}

// stopAnimation 停止动画
func (t *Grid) stopAnimation() {
	//当前格子的步数大于0的
	if 0 < t.movingCount {
		//将下一步直接赋值给当前
		t.current = t.next
		//下一步置为0
		t.next = GridData{}
	}
	//移动步数置为0
	t.movingCount = 0
	t.startPoppingCount = 0
	t.poppingCount = 0
}

// Move 移动
func (t *Grid) move() {
	t.movingCount--         //移动步数减1
	if t.movingCount == 0 { //移动到位置了
		if t.current.value != t.next.value && 0 < t.next.value { //判断当前的值是否等于移动后的值，并且移动后的值要大于0
			t.poppingCount = maxPoppingCount //将计数改为最大
		}
		t.current = t.next  //当前的格子更新为移动后的值
		t.next = GridData{} //移动后的置为0
	}
}

// Update 更新格子
func (t *Grid) Update() error {
	//循环，直到格子移动到边界
	switch {
	case 0 < t.movingCount: //格子移动步数大于0
		t.move()
	case 0 < t.startPoppingCount: //开始弹出计数大于0
		t.startPoppingCount-- //减1
	case 0 < t.poppingCount: //弹出计数大于0
		t.poppingCount-- //减1
	}
	return nil
}

// Draw draws the current tile to the given boardImage.
func (t *Grid) Draw(boardImage *ebiten.Image) {
	i, j := t.current.x, t.current.y
	ni, nj := t.next.x, t.next.y
	v := t.current.value
	if v == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin
	nx := ni*tileSize + (ni+1)*tileMargin
	ny := nj*tileSize + (nj+1)*tileMargin
	switch {
	case 0 < t.movingCount:
		rate := 1 - float64(t.movingCount)/maxMovingCount
		x = mean(x, nx, rate)
		y = mean(y, ny, rate)
	case 0 < t.startPoppingCount:
		rate := 1 - float64(t.startPoppingCount)/float64(maxPoppingCount)
		scale := meanF(0.0, 1.0, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	case 0 < t.poppingCount:
		const maxScale = 1.2
		rate := 0.0
		if maxPoppingCount*2/3 <= t.poppingCount {
			// 0 to 1
			rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount/3)
		} else {
			// 1 to 0
			rate = float64(t.poppingCount) / float64(maxPoppingCount*2/3)
		}
		scale := meanF(1.0, maxScale, rate)
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(gridBackgroundColor(v))
	boardImage.DrawImage(gridImage, op)
	str := strconv.Itoa(v)

	f := mplusBigFont
	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}

	w := font.MeasureString(f, str).Floor()
	h := (f.Metrics().Ascent + f.Metrics().Descent).Floor()
	x += (tileSize - w) / 2
	y += (tileSize-h)/2 + f.Metrics().Ascent.Floor()
	text.Draw(boardImage, str, f, x, y, gridColor(v))
}

func mean(a, b int, rate float64) int {
	return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	return a*(1-rate) + b*rate
}
