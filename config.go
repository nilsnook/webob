package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CONTENT_TYPE = "content_type"
	SERVER       = "server"
	USER_AGENT   = "user_agent"
	URL          = "url"
	STATUS_CODE  = "status_code"
	TICK         = "tick"
	defaultTick  = 60 * time.Second
)

type config struct {
	Url         string        `mapstructure:"url"`
	StatusCode  int           `mapstructure:"status_code"`
	ContentType string        `mapstructure:"content_type"`
	Server      string        `mapstructure:"server"`
	UserAgent   string        `mapstructure:"user_agent"`
	Tick        time.Duration `mapstructure:"tick"`
}

// func (c *config) setDefaults() {
// 	viper.SetDefault(CONTENT_TYPE, "")
// 	viper.SetDefault(SERVER, "")
// 	viper.SetDefault(USER_AGENT, "")
// 	viper.SetDefault(URL, "")
// 	viper.SetDefault(STATUS_CODE, 200)
// 	viper.SetDefault(TICK, defaultTick)
// }

func (c *config) readConfigFile() error {
	// set config name
	viper.SetConfigName("config")
	// set config type
	viper.SetConfigType("yaml")
	// set config path -
	// using user default config dir
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	appConfigDir := path.Join(userConfigDir, "webob")
	viper.AddConfigPath(appConfigDir)

	// read from config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *config) parseFlags() error {
	// define flags
	pflag.String(URL, "", "Request URL")
	pflag.String(STATUS_CODE, "", "HTTP response status code")
	pflag.String(CONTENT_TYPE, "", "Content-Type HTTP response header value")
	pflag.String(SERVER, "", "Server HTTP response header value")
	pflag.String(USER_AGENT, "", "User-Agent HTTP response header value")
	pflag.Duration(TICK, defaultTick, "Ticking interval")

	// parse flags
	pflag.Parse()
	// bind flags to viper
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}
	return nil
}

func (c *config) init() error {
	// read config file (if any)
	// ignoring error and continuing to parsing flags
	err := c.readConfigFile()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
	// parse flags
	err = c.parseFlags()
	if err != nil {
		return err
	}
	// unmarshal configuration to struct
	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *config) initFromConfigFile() error {
	// read from config file
	err := c.readConfigFile()
	if err != nil {
		return err
	}
	// unmarshal configuration to struct
	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}
