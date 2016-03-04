package interfacer
import (

  "github.com/jinzhu/gorm"
  "github.com/gorilla/sessions"
  //"github.com/gorilla/securecookie"
)



type AppContext struct {
  Db        *gorm.DB
  Store     *sessions.CookieStore
  //templates map[string]*template.Template
  //decoder *schema.Decoder
  // ... and the rest of our globals.
  //decoder   *schema.Decoder
  //bufpool   *bpool.Bufferpool
  
  //mandrill  *gochimp.MandrillAPI
  //log       *log.Logger
  //conf      *config // app-wide configuration: hostname, ports, etc.
}

