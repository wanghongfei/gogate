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
		fmt.Printf("path = %v, id = %s\n", servInfo.Prefix, servInfo.Id)

		return true
	})
}

func TestRouter_Match(t *testing.T) {
	r, _ := NewRouter("../route.yml")

	result := r.Match("/user")
	fmt.Println(result)
	if "user-service" != result {
		t.Errorf("/user mismatch, %s\n", result)
	}

	result = r.Match("/order")
	fmt.Println(result)
	if "order-service" != result {
		t.Errorf("/order mismatch, %s\n", result)
	}

	result = r.Match("/user/dog")
	fmt.Println(result)

	result = r.Match("/nomatch")
	fmt.Println(result)

	result = r.Match("/us")
	fmt.Println(result)
}
