package utils

import "github.com/bwmarrin/snowflake"

var generator *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if nil != err {
		panic(err)
	}

	generator = node
}

func GenerateUuid() int64 {
	return generator.Generate().Int64()
}

