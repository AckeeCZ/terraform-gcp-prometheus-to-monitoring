package main

import (
	"fmt"
	"os"
)

func load_var(config map[string]*string, name string) {
	value, exists := os.LookupEnv(name)
	if exists {
		config[name] = &value
	} else {
		panic(fmt.Sprintf("%s not present", name))
	}
}

func load_vars(config map[string]*string) {
	load_var(config, "PROJECT_ID")
	load_var(config, "SERVICE")
}
