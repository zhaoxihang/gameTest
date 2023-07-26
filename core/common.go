package core

// abs 绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// vecToDir  向那个方向移动，是否移动
func vecToDir(dx, dy int) (Dir, bool) {
	//格子是4*4的 0123*0123
	//移动的位置都要小于格子对应的边界
	if abs(dx) < boardSize && abs(dy) < boardSize {
		return 0, false //两个值的绝对值都小于
	}
	//
	if abs(dx) < abs(dy) { //第一个值的绝对值小于第二个值的绝对值
		if dy < 0 {
			return DirUp, true //第二个值小于0
		}
		return DirDown, true //第二个值大于0
	}
	if dx < 0 { //第一个值小于0
		return DirLeft, true
	}
	return DirRight, true
}
