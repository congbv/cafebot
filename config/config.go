package config

import (
	"os"

	congbvlog "log"

	"github.com/jessevdk/go-flags"
)

// Config holds all application configuration data
type Config struct {
	Port       string `short:"p" long:"port" description:"port to listen on" default:"8080"`
	TgAPIToken string `long:"tg-api-token" description:"telegram bot api token" default:"981238424:AAFXiUGS1yQqqSiHMjnnSfPg_mWFXON7gtA"`
	LogLevel   string `long:"log-level" description:"log level for all subsystems {trace, debug, info, error, critical}" default:"info"`

	CafeConfigFile string `long:"cafe-config" description:"cafe config file path" default:"cafe.json"`

	Cafe CafeConfig
}

//Load ok
func Load() (Config, error) {
	conf := Config{}
	_, err := flags.Parse(&conf)

	congbvlog.Printf("conf in load %+v \n", conf)

	// go-flags package should consider error wrapping instead of this
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		os.Exit(0)
	}
	conf.Cafe, err = loadCafeConfig(conf.CafeConfigFile)
	return conf, err
}
