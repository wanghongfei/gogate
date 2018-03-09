package utils

import (
	"fmt"
	"testing"
)

func TestGetFirstNoneLoopIp(t *testing.T) {
	ip, err := GetFirstNoneLoopIp()
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(ip)
}
