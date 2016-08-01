package interfacer

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/robfig/config"
	log "gopkg.in/inconshreveable/log15.v2"
)

// AppConfig holds and manages configurations of this app
type AppConfig struct {
	cf     *conf
	loader *config.Config
}

type conf struct {
	Mode     string
	BindIP   string
	BindPort string
	Language string
	Secret   string
}

// LoadConfig loads application configuration from file and return Config object
func LoadConfig() *AppConfig {
	loader, err := config.ReadDefault("./ladybug.conf")
	if err != nil {
		log.Error("conf", "msg", err.Error())
	}

	var appConf AppConfig
	var conf conf

	modePtr := flag.String("mode", "dev", "a string")

	flag.Parse()

	conf.Mode = *modePtr
	log.Info("Config", "Starting Mode", conf.Mode)

	if rv, err := loader.String("", "app.secret"); err != nil {
		log.Error("conf", "msg", "fail to load app secret")
	} else {
		conf.Secret = rv
	}

	appConf.cf = &conf
	appConf.loader = loader

	return &appConf
}

// GetMode returns mode of the app(ex : dev, prod ...)
func (conf AppConfig) GetMode() string {
	if conf.cf == nil {
		return ""
	}
	return conf.cf.Mode
}

// GetBindAddress returns bind address of this app
// default value is localhost:8080 if the value is not set
// For example : dev.wisedog.net:80
func (conf AppConfig) GetBindAddress() string {
	if conf.cf == nil {
		return "localhost:8080"
	}
	var bindIP string
	var bindPort string
	if rv, err := conf.loader.String(conf.cf.Mode, "http.addr"); err != nil {
		log.Error("conf", "msg", "fail to load address")
	} else {
		bindIP = rv
	}

	if rv, err := conf.loader.String(conf.cf.Mode, "http.port"); err != nil {
		log.Error("conf", "msg", "fail to load port")
	} else {
		bindPort = rv
	}

	// heroku support
	if strings.Contains(conf.GetMode(), "heroku") {
		bindIP = ""
		bindPort = os.Getenv("PORT")
		if bindPort == "" {
			log.Error("Config", "msg", "Fail to get PORT environment")
		}
	}

	s := fmt.Sprintf("%s:%s", bindIP, bindPort)
	return s
}

// GetValue returns the appropriate value with given key
// if any value matched the given key, returns empty string
func (conf AppConfig) GetValue(key string) string {
	if rv, err := conf.loader.String(conf.cf.Mode, key); err != nil {
		log.Error("conf", "msg", "fail to load port")
	} else {
		return rv
	}
	return ""
}
