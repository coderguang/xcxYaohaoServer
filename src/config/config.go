package yaohaoConfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	yaohaoDef "xcxYaohaoServer/src/define"
)

func ReadConfig(configfile string) *yaohaoDef.Config {
	config, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Println("read config error")
		os.Exit(1)
	}
	t := new(yaohaoDef.Config)
	p := &t
	err = json.Unmarshal([]byte(config), p)
	if err != nil {
		log.Println("parse config error")
		os.Exit(1)
	}
	return t
}
