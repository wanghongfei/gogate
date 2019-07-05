package route

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadRoute(t *testing.T) {
	//routeMap, _, err := loadRoute("../route.yml")
	//if nil != err {
	//	t.Error(err)
	//}
	//
	//for _, servInfo := range routeMap {
	//	fmt.Printf("path = %v, id = %s\n", servInfo.Prefix, servInfo.Id)
	//}
}

func TestRouter_Match(t *testing.T) {
	r, err := NewRouter("../../route.yml")
	if nil != err {
		t.Fatal(err)
	}

	result := r.Match("/user")
	fmt.Println(result)
	if "user-service" != result.Id {
		t.Errorf("/user mismatch, %s\n", result)
	}

	result = r.Match("/order")
	fmt.Println(result)
	if "order-service" != result.Id {
		t.Errorf("/order mismatch, %s\n", result)
	}

	result = r.Match("/aaaa")
	if nil != result {
		t.Errorf("/aaaa mismatch, %s\n", result)
	}
	fmt.Println(result)

	result = r.Match("/img")
	if "localhost:8080,localhost:8080" != result.Host {
		t.Errorf("/img mismatch, %s\n", result)
	}
	fmt.Println(result)
}

func BenchmarkRouter_Match(b *testing.B) {
	r, err := NewRouter("../../route.yml")
	if nil != err {
		b.Fatal(err)
	}


	for ix := 0; ix < b.N; ix++ {
		r.Match("/order/a/b/c/d/e/f/g")
	}
}

func TestYaml(t *testing.T) {
	f, err := os.Open("../route.yml")
	if nil != err {
		t.Error(err)
		return
	}
	defer f.Close()

	buf, _ := ioutil.ReadAll(f)

	yamlMap := make(map[string]interface{})
	yaml.Unmarshal(buf, &yamlMap)

	fmt.Println(yamlMap["services"])
}
