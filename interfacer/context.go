package interfacer
import (

  "github.com/jinzhu/gorm"
  "github.com/gorilla/sessions"
  "github.com/gorilla/schema"
  //"github.com/gorilla/securecookie"
  "github.com/wisedog/ladybug/models"
)



type AppContext struct {
  Db        *gorm.DB
  Store     *sessions.CookieStore
  Decoder   *schema.Decoder

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