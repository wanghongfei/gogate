package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/alecthomas/log4go"
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

func LoadConfig(filename string) {
	f, err := os.Open(filename)
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

func InitLog(filename string) {
	log4go.LoadConfiguration(filename)
}

