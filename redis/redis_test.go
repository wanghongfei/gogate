package redis

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

)

func TestGetString(t *testing.T) {
	c := NewRedisClient("127.0.0.1:6379", 1)
	err := c.Connect()
	if nil != err {
		t.Error(err)
		return
	}


	str, err := c.GetString("abc")
	fmt.Println(str)
	c.Close()
}

func TestRedisClient_ExeLuaInt(t *testing.T) {
	c := NewRedisClient("127.0.0.1:6379", 1)
	err := c.Connect()
	if nil != err {
		t.Error(err)
		return
	}
	defer c.Close()

	luaFile, err := os.Open("../lua/rate_limiter.lua")
	if nil != err {
		t.Error(err)
		return
	}
	defer luaFile.Close()

	luaBuf, _ := ioutil.ReadAll(luaFile)
	resp, err := c.ExeLuaInt(string(luaBuf), nil, []string{"10"})
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(resp)
}
