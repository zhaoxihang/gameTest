package core

// Dir represents a direction.
type Dir int //方向

const (
	DirUp    Dir = iota //上
	DirRight            //右
	DirDown             //下
	DirLeft             //左
)

// String returns a string representing the direction.
func (d Dir) String() string {
	switch d {
	case DirUp:
		return "Up"
	case DirRight:
		return "Right"
	case DirDown:
		return "Down"
	case DirLeft:
		return "Left"
	}
	panic("not reach")
}

// Vector returns a [-1, 1] value for each axis.
func (d Dir) Vector() (x, y int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirRight:
		return 1, 0
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	}
	panic("not reach")
}
