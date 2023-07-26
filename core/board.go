package core

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"sort"
)

var taskTerminated = errors.New("twenty48: task terminated")

type task func() error

// Board 游戏棋盘
type Board struct {
	size  int //棋盘大小
	grids map[*Grid]struct{}
	tasks []task
}

//  0  1  2  3
//  4  5  6  7
//  8  9 10 11
// 12 13 14 15

// NewBoard 初始化棋盘
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size:  size,
		grids: map[*Grid]struct{}{},
	}
	//第一次增加两个格子
	for i := 0; i < 2; i++ {
		//给棋盘增加格子
		if err := b.addRandomGrid(); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// addRandomGrid 增加随机的格子
func (b *Board) addRandomGrid() error {
	//初始化一个棋盘的位置
	cells := make([]bool, b.size*b.size)
	for grid := range b.grids {
		//判断已有的格子中是否存在有步数的格子
		if grid.IsMoving() {
			panic("not reach")
		}
		//计算格子的位置
		i := grid.current.x + grid.current.y*b.size
		//棋盘中该位置有值
		cells[i] = true
	}
	//初始化一个没有棋子的空格子集
	var availableCells []int
	for i, b := range cells {
		//格子有值就跳过
		if b {
			continue
		}
		//空值的写入集合
		availableCells = append(availableCells, i)
	}
	//判断是否还有空格子
	if len(availableCells) == 0 {
		return errors.New("twenty48: there is no space to add a new tile")
	}
	//随机取出一个位置
	c := availableCells[rand.Intn(len(availableCells))]
	//格子的值为2
	v := 2
	// 1/10 的概率为4
	if rand.Intn(10) == 0 {
		v = 4
	}
	// 计算格子的x,y轴
	x := c % b.size
	y := c / b.size
	// 初始化格子
	t := NewGrid(v, x, y)
	// 写入棋盘
	b.grids[t] = struct{}{}
	return nil
}

// Update 更新棋盘状态
func (b *Board) Update(input *Input) error {
	for t := range b.grids {
		//更新格子状态
		if err := t.Update(); err != nil {
			return err
		}
	}
	//判断是否有任务
	if 0 < len(b.tasks) {
		//取出第一个任务并执行
		t := b.tasks[0]
		if err := t(); err == taskTerminated {
			//范围值为已终止的将任务删除
			b.tasks = b.tasks[1:]
		} else if err != nil {
			return err
		}
		return nil
	}
	//计算输入的移动
	if dir, ok := input.Dir(); ok {
		//棋盘开始移动
		if err := b.Move(dir); err != nil {
			return err
		}
	}
	return nil
}

// Move 将棋盘的移动入队
func (b *Board) Move(dir Dir) error {
	for t := range b.grids {
		t.stopAnimation()
	}
	//移动格子
	if !b.MoveGrids(dir) {
		return nil
	}
	//移动成功
	b.tasks = append(b.tasks, func() error {
		//将每个格子判断是否需要移动的写入任务
		for t := range b.grids {
			if t.IsMoving() {
				return nil
			}
		}
		return taskTerminated
	})
	b.tasks = append(b.tasks, func() error {
		nextTiles := map[*Grid]struct{}{}
		for t := range b.grids {
			//格子需要移动的报错
			if t.IsMoving() {
				panic("not reach")
			}
			//格子的下一次移动的值不等于0的报错
			if t.next.value != 0 {
				panic("not reach")
			}
			//格子的值等于0的跳过
			if t.current.value == 0 {
				continue
			}
			//格子不动的写入
			nextTiles[t] = struct{}{}
		}
		//更新格子
		b.grids = nextTiles
		//增加随机的格子
		if err := b.addRandomGrid(); err != nil {
			return err
		}
		return taskTerminated
	})
	return nil
}

// gridAt 找到该位置的格子
func (b *Board) gridAt(x, y int) *Grid {
	var result *Grid
	for t := range b.grids {
		//如果当前格子的x轴不等于x，或者y轴不等于y，跳过
		if t.current.x != x || t.current.y != y {
			continue
		}
		//格子有值，报错
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

// currentOrNextGridAt 找到移动后的格子的位置
func (b *Board) currentOrNextGridAt(x, y int) *Grid {
	var result *Grid
	for t := range b.grids {
		//移动的步数是否大于0
		if t.IsMoving() {
			//判断移动后位置的 x轴不等于x的，y轴不等于y的，值等于0的跳过
			if t.next.x != x || t.next.y != y || t.next.value == 0 {
				continue
			}
		} else {
			//判断当前位置的 x轴不等于x的，y轴不等于y的跳过
			if t.current.x != x || t.current.y != y {
				continue
			}
		}
		//格子有值的报错
		if result != nil {
			panic("not reach")
		}
		result = t
	}
	return result
}

// MoveGrids 移动格子集合 返回移动是否成功
func (b *Board) MoveGrids(dir Dir) bool {
	tiles := b.grids
	size := b.size
	//移动格子的矢量坐标
	vx, vy := dir.Vector()
	var tx []int
	var ty []int
	//计算棋盘每个格子的坐标
	for i := 0; i < size; i++ {
		tx = append(tx, i)
		ty = append(ty, i)
	}
	//x轴是否需要移动
	if vx > 0 {
		//对x轴逆序排序
		sort.Sort(sort.Reverse(sort.IntSlice(tx)))
	}
	//y轴是否需要移动
	if vy > 0 {
		//对y轴逆序排序
		sort.Sort(sort.Reverse(sort.IntSlice(ty)))
	}
	//定义一个标志位：是否需要移动
	moved := false
	//对y轴进行循环
	for _, j := range ty {
		//对x轴进行循环
		for _, i := range tx {
			//找到该位置的格子
			t := b.gridAt(i, j)
			if t == nil { //格子为空跳出循环
				continue
			}
			//如果移动后的格子不为空 报错
			if t.next != (GridData{}) {
				panic("not reach")
			}
			//如果移动步数不为0 报错
			if t.IsMoving() {
				panic("not reach")
			}
			// (ii, jj) 是格子t的下一个位置。
			// (ii, jj) 被更新，直到找到可合并切片或格子t不能再移动了。
			ii := i
			jj := j
			for {
				//计算移动后的位置
				ni := ii + vx
				nj := jj + vy
				//移动后的位置不能超过框
				if ni < 0 || ni >= size || nj < 0 || nj >= size {
					break
				}
				//找到移动后的格子的位置
				tt := b.currentOrNextGridAt(ni, nj)
				if tt == nil { //格子等于空的
					//格子移动后的位置
					ii = ni
					jj = nj
					//标志，需要位移
					moved = true
					continue
				}
				//当前格子的值不等于移动后的位置的值，跳过
				if t.current.value != tt.current.value {
					break
				}
				//移动后位置的移动值大于0的，并且当前值不等于移动后值的跳过
				if 0 < tt.movingCount && tt.current.value != tt.next.value {
					// tt is already being merged with another tile.
					// Break here without updating (ii, jj).
					break
				}
				ii = ni
				jj = nj
				moved = true
				break
			}
			// 下一步是格子t的下一状态。
			next := GridData{}
			next.value = t.current.value
			// 如果下一个位置(II，JJ)有格子，则应为可合并。让我们合并吧。
			if tt := b.currentOrNextGridAt(ii, jj); tt != t && tt != nil {
				next.value = t.current.value + tt.current.value
				tt.next.value = 0
				tt.next.x = ii
				tt.next.y = jj
				tt.movingCount = maxMovingCount
			}
			next.x = ii
			next.y = jj
			if t.current != next {
				t.next = next
				t.movingCount = maxMovingCount
			}
		}
	}
	if !moved {
		for t := range tiles {
			t.next = GridData{}
			t.movingCount = 0
		}
	}
	return moved
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}

func (b *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			v := 0
			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(gridBackgroundColor(v))
			boardImage.DrawImage(gridImage, op)
		}
	}
	animatingTiles := map[*Grid]struct{}{}
	nonAnimatingTiles := map[*Grid]struct{}{}
	for t := range b.grids {
		if t.IsMoving() {
			animatingTiles[t] = struct{}{}
		} else {
			nonAnimatingTiles[t] = struct{}{}
		}
	}
	for t := range nonAnimatingTiles {
		t.Draw(boardImage)
	}
	for t := range animatingTiles {
		t.Draw(boardImage)
	}
}
