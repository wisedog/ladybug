package interfacer
import (

  "github.com/jinzhu/gorm"
  "github.com/gorilla/sessions"
  "github.com/gorilla/schema"
  "github.com/wisedog/ladybug/models"
  //"github.com/robfig/config"
)



type AppContext struct {
  Db        *gorm.DB
  Store     *sessions.CookieStore
  Decoder   *schema.Decoder
  Config    *AppConfig

  // project-specific fields
  User      *models.User
  ProjectName string
  //templates map[string]*template.Template
  // ... and the rest of our globals.
  //bufpool   *bpool.Bufferpool
  
  //mandrill  *gochimp.MandrillAPI
  //log       *log.Logger
  //conf      *config // app-wide configuration: hostname, ports, etc.
}