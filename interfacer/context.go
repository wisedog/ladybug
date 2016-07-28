package interfacer

import (
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/wisedog/ladybug/models"
	//"github.com/robfig/config"
)

// AppContext is a context of each request
type AppContext struct {
	Db      *gorm.DB
	Store   *sessions.CookieStore
	Decoder *schema.Decoder
	Config  *AppConfig

	// project-specific fields
	User        *models.User
	ProjectName string
	//templates map[string]*template.Template
	// ... and the rest of our globals.
	//bufpool   *bpool.Bufferpool

	//mandrill  *gochimp.MandrillAPI
	//log       *log.Logger
	//conf      *config // app-wide configuration: hostname, ports, etc.
}
