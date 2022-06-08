package p

import (
	"fmt"
	"os"

	zlog "github.com/rs/zerolog/log"
)

func load_var(config map[string]*string, name string) {
	value, exists := os.LookupEnv(name)
	if exists {
		config[name] = &value
	} else {
		zlog.Print(fmt.Sprintf("%s not present", name))
		os.Exit(1)
	}
}

func load_vars(config map[string]*string) {
	load_var(config, "PROJECT_ID")
}
