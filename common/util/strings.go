package util

import "strings"

// StringConcat 使用strings.Builder预分配进行字符串连接
func StringConcat(arg ...string) string {
	var sb strings.Builder

	if len(arg) == 0 {
		return ""
	}

	// 预分配
	sb.Grow(len(arg))

	for i := 0; i < len(arg); i++ {
		sb.WriteString(arg[i])
	}

	return sb.String()
}
