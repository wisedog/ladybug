package interfacer

import (
	"fmt"

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
	BindPort int
	Dialect  string
	Language string
	Secret   string
}

// LoadConfigWithArgs initializes with addr, port, db dialect
// and other options with application configuration from file and return Config object
// if mode is empty, load only default options (root section)
// mode is a string indicates section of config file.
// filename is config file, especially used by unit testing
func LoadConfigWithArgs(mode, addr string, port int, dialect, filename string) *AppConfig {
	loader, err := config.ReadDefault(filename)
	if err != nil {
		log.Error("conf", "msg", err.Error())
	}

	var appConf AppConfig
	var conf conf

	conf.Mode = mode

	if rv, err := loader.String("", "app.secret"); err != nil {
		log.Error("conf", "msg", "fail to load app secret")
	} else {
		conf.Secret = rv
	}
	// filtering default values
	conf.BindIP = addr
	conf.BindPort = port
	conf.Dialect = dialect

	appConf.cf = &conf
	appConf.loader = loader

	return &appConf
}

// LoadConfig loads application configuration from file and return Config object
// mode is a string indicates section of config file.
// filename is config file, used by unit testing
// Binding address, port are initialized by localhost, 8000
func LoadConfig(mode, filename string) *AppConfig {
	loader, err := config.ReadDefault(filename)
	if err != nil {
		log.Error("conf", "msg", err.Error())
	}

	var appConf AppConfig
	var conf conf

	conf.Mode = mode

	if rv, err := loader.String("", "app.secret"); err != nil {
		log.Error("conf", "msg", "fail to load app secret")
	} else {
		conf.Secret = rv
	}

	// filtering default values
	conf.BindIP = "localhost"
	conf.BindPort = 8000
	conf.Dialect = ""

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

	s := fmt.Sprintf("%s:%d", conf.cf.BindIP, conf.cf.BindPort)
	return s
}

// GetValue returns the appropriate value with given key
// if any value matched the given key, returns empty string
func (conf AppConfig) GetValue(key string) string {
	if rv, err := conf.loader.String(conf.cf.Mode, key); err != nil {
		log.Error("conf", "msg", "fail to load key", "key", key, "mode", conf.cf.Mode)
	} else {
		return rv
	}
	return ""
}

// GetDialect returns database dialect
func (conf AppConfig) GetDialect() string {
	return conf.cf.Dialect
}
