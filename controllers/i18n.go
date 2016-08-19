package controllers

import (
	"github.com/robfig/config"
	log "gopkg.in/inconshreveable/log15.v2"
)

// LoadI18nMessage loads i18n resources. i18n messages is at /messages
// TODO support other languages
func LoadI18nMessage() *config.Config {
	c, err := config.ReadDefault("messages/ladybug.en")
	if err != nil {
		log.Error("i18n", "msg", err.Error())
	}

	return c
}

// GetI18nMessage returns i18n string correspond to given key
func GetI18nMessage(key string) string {
	if messages == nil {
		messages = LoadI18nMessage()
	}
	s, _ := messages.String("", key)
	return s
}
