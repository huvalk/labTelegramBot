package config

import (
	"github.com/kataras/golog"
	"io/ioutil"
)

func GetDBConfig() string {
	content, err := ioutil.ReadFile("config/initdb.sql")
	if err != nil {
		golog.Fatal("Error with config: ", err)
	}

	return string(content)
}
