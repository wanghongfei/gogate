package utils


import (
	"fmt"
	"runtime"
	"strconv"
)

// 创建一个携带环境信息的error, 包含文件名 + 行号
func Errorf(format string, args ...interface{}) error {
	_, srcName, line, _ := runtime.Caller(1)

	// [文件名:行号]
	prefix := "[" + srcName + ":" + strconv.Itoa(line) + "] "
	return fmt.Errorf(prefix + format, args...)
}

