package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// for test only
func Connect() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}

	defer c.Close()
}
