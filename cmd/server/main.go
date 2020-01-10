package main

import (
	"FinalTask/Internal/App/Apiserver"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/Apiserver.toml", "path to config file") //парсим configPath,флаг "config-path",значение по умолчанию,описание для хелпа
}

func main() {
	flag.Parse()

	config := Apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
	s := Apiserver.New(config)
	fmt.Println(s)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
