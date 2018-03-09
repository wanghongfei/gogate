package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"code.google.com/p/log4go"
)

type AppConfig struct {
	AppName			string
	Host			string
	Port			int
	MaxConnection	int

	EurekaConfig	string
	RouteConfig		string
}

var App *AppConfig

func init() {
	f, err := os.Open("gogate.json")
	if nil != err {
		log4go.Error(err)
		os.Exit(1)
	}
	defer f.Close()

	buf, _ := ioutil.ReadAll(f)
	App = new(AppConfig)
	err = json.Unmarshal(buf, App)

	if nil != err {
		log4go.Error(err)
		os.Exit(1)
	}
}
