package util

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

/*********************************************************
** 函数功能: 比较字符串数值大小
** 参    数：两个数值型字符串
** 返 回 值: 前者是否大于后者
**********************************************************/
func StringCompare(a, b string) bool {
	if len(a) == len(a) {
		if a > b {
			return true
		} else {
			return false
		}
	} else {
		if len(a) > len(b) {
			return true
		} else {
			return false
		}
	}
}
