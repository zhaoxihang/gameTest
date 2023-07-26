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
	tileSize   = 80 //每个格子的宽高
	tileMargin = 4  //每个格子之间的间距
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

	movingCount       int //移动的步数 动画展示几次到最终目的地，动画几帧移动到目的地
	startPoppingCount int //弹出数字经过几帧
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

// Draw 将当前格子绘制到给定的boardImage。
func (t *Grid) Draw(boardImage *ebiten.Image) {
	//获取当前格子的位置
	i, j := t.current.x, t.current.y
	//获取移动后的位置
	ni, nj := t.next.x, t.next.y
	//获取当前的值
	v := t.current.value
	//当前值等于0，不更新
	if v == 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin    //计算当前格子的x轴左边位置
	y := j*tileSize + (j+1)*tileMargin    //计算当前格子的y轴上边位置
	nx := ni*tileSize + (ni+1)*tileMargin //计算移动后格子的x轴左边位置
	ny := nj*tileSize + (nj+1)*tileMargin //计算移动后格子的y轴上边位置
	switch {
	case 0 < t.movingCount: //移动
		//每次向前移动1/5
		rate := 1 - float64(t.movingCount)/maxMovingCount
		x = mean(x, nx, rate)
		y = mean(y, ny, rate)
	case 0 < t.startPoppingCount: //生成
		rate := 1 - float64(t.startPoppingCount)/float64(maxPoppingCount)
		scale := meanF(0.0, 1.0, rate)
		//格子慢慢变大
		op.GeoM.Translate(float64(-tileSize/2), float64(-tileSize/2))
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(tileSize/2), float64(tileSize/2))
	case 0 < t.poppingCount: //合并
		//合并的时候变大一下
		const maxScale = 1.2
		rate := 0.0
		if maxPoppingCount*2/3 <= t.poppingCount {
			// 0 to 1 前两帧
			// 1-(6-2*6/3)/(6/3)=1-2/4=0.5
			// 1-(5-2*6/3)/(6/3)=1-1/4=0.75
			rate = 1 - float64(t.poppingCount-2*maxPoppingCount/3)/float64(maxPoppingCount/3)
		} else {
			// 1 to 0 后四帧
			// 4/(6*2/3)=4/4=1
			// 3/(6*2/3)=3/4=0.75
			// 2/(6*2/3)=2/4=0.5
			// 1/(6*2/3)=1/4=0.25
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
	//格子中的值转换为字符串
	str := strconv.Itoa(v)

	f := mplusBigFont
	//值长度超过2用普通字体
	//值长度超过3用小字体
	//其余使用大字体
	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}
	//计算字体的位置
	w := font.MeasureString(f, str).Floor()
	h := (f.Metrics().Ascent + f.Metrics().Descent).Floor()
	//居中
	x += (tileSize - w) / 2
	y += (tileSize-h)/2 + f.Metrics().Ascent.Floor()
	text.Draw(boardImage, str, f, x, y, gridColor(v))
}

// mean 计算a移动到b,走过rate后的值
func mean(a, b int, rate float64) int {
	//a-a*rate+b*rate=a-(a+b)*rate                 | b-(b-a)*(1-rate) = b-(b-b*rate-a+a*rate) = b-b+b*rate+a-a*rate = b*rate+a-a*rate = a+(b-a)*rate
	//假设 a= 10 b=20                     		   | a+(b-a)*rate
	//1: 10*(1-0.2)+20*0.2 = 8 + 4 = 12   ->2      |  10+(20-10)*0.2=12
	//2: 12*(1-0.4)+20*0.4 = 7.2 + 8 = 15 ->3      |  12+(20-12)*0.4=15
	//3: 15*(1-0.6)+20*0.6 = 6 + 12 = 18 ->3       |  15+(20-15)*0.6=18
	//4: 18*(1-0.8)+20*0.8 = 3.6 + 16 = 19 ->1     |  18+(20-18)*0.8=19
	//5: 19*(1-1)+20*1 = 0 + 20 = 20 ->1      	   |  19+(20-19)*1=20

	return int(float64(a) + float64(b-a)*rate)
	//return int(float64(a)*(1-rate) + float64(b)*rate)
}

func meanF(a, b float64, rate float64) float64 {
	//放大a到b (b-a)*rate的倍数
	//1*(1-rate) + 1.2*rate = 1-rate + 1.2*rate = 1 + 0.2*rate
	return a + (b-a)*rate
	//return a*(1-rate) + b*rate
}
