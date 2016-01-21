package main

import (
	"flag"
	"fmt"
	"jd.com/cdn/cashgame/cfg"
	// "jd.com/cdn/cashgame/model"
	"jd.com/cdn/cashgame/web"
)

type flags struct {
	config *string
}

func parseFlags() *flags {
	flags := &flags{}
	flags.config = flag.String("conf", "./config.json", "config file")
	flag.Parse()

	return flags
}

func initConfig(configPath string) (*cfg.Config, error) {
	return cfg.ParseConfig(configPath)
}

func main() {
	flags := parseFlags()
	fmt.Printf("config file: %s\n", *flags.config)

	conf, err := initConfig(*flags.config)
	if err != nil {
		fmt.Printf("error init config: %v\n", err)
		return
	}
	fmt.Printf("succeed init config\n")

	web.InitEnv(conf)

	err = web.WebStart(conf)
	if err != nil {
		fmt.Printf("WebServer running failed, error: %v \n", err)
	}
}
