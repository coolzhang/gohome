package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type tomlConfig struct {
	Ordercenter dsn
	Runaway     dsn
}

type dsn struct {
	User     string
	Password string
	Host     string
	Port     string
}

func tomlParse() {
	var conf tomlConfig
	if _, err := toml.DecodeFile("dsn.toml", &conf); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s:%s@tcp(%s:%s)\n", conf.Ordercenter.User, conf.Ordercenter.Password, conf.Ordercenter.Host, conf.Ordercenter.Port)
	fmt.Printf("%s:%s@tcp(%s:%s)\n", conf.Runaway.User, conf.Runaway.Password, conf.Runaway.Host, conf.Runaway.Port)
}

func viperParse() {
	viper.SetConfigType("toml")
	viper.SetConfigName("dsn")
	viper.AddConfigPath("/data/mygo/src/gohome/myweb")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Printf("%s:%s@tcp(%s:%s)\n", viper.GetStringMapString("runaway")["user"], viper.GetStringMapString("runaway")["password"], viper.GetStringMapString("runaway")["host"], viper.GetStringMapString("runaway")["port"])
	fmt.Printf("%s:%s@tcp(%s:%s)\n", viper.GetStringMapString("ordercenter")["user"], viper.GetStringMapString("ordercenter")["password"], viper.GetStringMapString("ordercenter")["host"], viper.GetStringMapString("ordercenter")["port"])
}

func main() {
	tomlParse()
	//viperParse()
}
