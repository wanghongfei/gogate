package serv

import (
	"fmt"
	"testing"
)

func TestLoadRoute(t *testing.T) {
	routeMap, err := loadRoute("../route.yml")
	if nil != err {
		t.Error(err)
	}

	routeMap.Range(func(name, info interface{}) bool {
		servInfo := info.(*ServiceInfo)
		fmt.Printf("path = %v, id = %s\n", servInfo.Path, servInfo.Id)

		return true
	})
}
