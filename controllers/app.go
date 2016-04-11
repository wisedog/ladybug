package controllers

import (
  "fmt"
  "os"
  "strings"

  "net/http"
  "html/template"
  "path/filepath"
  "encoding/gob"

  "golang.org/x/crypto/bcrypt"
  "github.com/robfig/config"
  "github.com/gorilla/sessions"
  "github.com/wisedog/ladybug/errors"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  
  log "gopkg.in/inconshreveable/log15.v2"
)

const (
  LadybugSession = "ladybug_session"
  FlashAuthMessage = "flash_auth_msg"
)


// function map used in template files
var(
  funcMap = template.FuncMap {
    "nl2br": func (str string) string { return strings.Replace(str, "\n", "<br />", -1) },
    "msg" : func (key string) string{
       if messages == nil{
        messages = LoadI18nMessage()
      }
      s, _ := messages.String("", key)
      return s
    },
    "categoryi18n" : func(key string) string{
        return ""
    },  
  }
  messages *config.Config
)

func init(){

  // for flash messages
  gob.Register(&models.Build{})
  gob.Register(&models.Project{})
  gob.Register(&map[string]string{})

  // TODO store category map to static variable. 
  // in funcMap - categorui18n function refers the var.
}

//@deprecated
//var store = sessions.NewCookieStore([]byte("something-very-secret"))

/*var cookieHandler = securecookie.New(
  securecookie.GenerateRandomKey(64),
  securecookie.GenerateRandomKey(32),
  )*/

// connected is private utilitiy function for checking 
// this user id now on connected or not.
func connected(c *interfacer.AppContext, r *http.Request) *models.User{
  session, err := c.Store.Get(r, "ladybug")
  if err != nil {
    log.Warn("An error in connected(),", "msg", err)
    return nil
  }
  email := session.Values["user"]
  uid := session.Values["uid"]

  var u models.User
  c.Db.Where("email = ? and id = ?", email, uid ).First(&u)

  if u.ID == 0{
    return nil
  }
  return &u
}

func renderTemplate(w http.ResponseWriter, tmpl string, msg string, isLogout bool ) {
  dir, _ := os.Getwd()
  target := filepath.Join(dir, "views", "application", tmpl )
  t, err := template.ParseFiles(target)

  if err != nil{
    log.Error("App", "msg", "Template parse files error", "err", err)
  }

  items := map[string]interface{}{
    "Messages": msg,
    "IsLogout" : isLogout,
  }

  fmt.Println("items", items)
  
  if err := t.Execute(w, items); err != nil{
    log.Error("App", "msg", "Template Execute error","err", err)
  }
}

// Login checks user email, password 
func Login(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

  session, err := c.Store.Get(r, "ladybug")
  if err != nil {
      http.Error(w, err.Error(), 500)
      return nil
  }

  email := r.FormValue("email")
  passwd := r.FormValue("password")

  var user models.User
  
  if err := c.Db.Where("email = ?", email).First(&user); err != nil {
    
    if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(passwd)); err == nil {
      session.Values["user"] = email
      session.Values["uid"] = user.ID

      fmt.Println("user", user)
      fmt.Println("uid", user.ID)

      // Save it before we write to the response/return from the handler.
      session.Save(r, w)

      http.Redirect(w, r, "/hello", http.StatusFound)
    }
  }

  // not exist account or incorrect password
  session.AddFlash("Not exist account or incorrect password", FlashAuthMessage)
  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
  return nil
}

// LoginPage renders a login page
func LoginPage(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  session, err := c.Store.Get(r, "ladybug")
  if err != nil {
      http.Error(w, err.Error(), 500)
      return err
  }

  // if already logged in, redirect to Hello
  if session.Values["user"] != nil && session.Values["uid"] != nil {
    fmt.Println("Already logged in. Loginpage user/uid", session.Values["user"], session.Values["uid"])
    email := session.Values["user"]
    uid := session.Values["uid"]

    var u models.User
    c.Db.Where("email = ? and id = ?", email, uid ).First(&u)

    if u.ID == 0{
      return nil  //TODO return unauth
    }

    http.Redirect(w, r, "/hello", http.StatusFound)

  }else{
    log.Debug("App", "msg", "not login, rendering login page")

    // Get the previously flashes, if any.
    var s string
    var logoutFlag bool

    if msg := session.Flashes(FlashAuthMessage); len(msg) > 0{
      v, ok := msg[0].(string)
      if ok {
        s = v
      }
    }
    renderTemplate(w, "index.html", s, logoutFlag)
  }
  
  return nil
}

// LogOut func clears login session, redirect to login page
func LogOut(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  session, err := c.Store.Get(r, "ladybug")
  if err != nil {
      http.Error(w, err.Error(), 500)
      return err
  }
  session.Values["user"] = ""
  session.Values["uid"] = 0
  session.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   -1,
  }
  log.Debug("App", "msg", "logout")
  session.AddFlash("You just logged out.", FlashAuthMessage)

  session.Save(r, w)
  http.Redirect(w, r, "/", http.StatusFound)
  return nil
}

// LogAndHTTPError makes log messages and return http error with given status http code
// default log level is ERROR
func LogAndHTTPError(status int, module string, errType string, msg string) error{
  log.Error(module, "type" , errType, "msg", msg)
  return errors.HttpError{Status : status, Desc : msg}
}


/*
func (c Application) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Application.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	c.Tx.NewRecord(user)
	c.Tx.Create(user)

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.Users.Index(user.ID))
}

func (c Application) Login(email, password string, remember bool) revel.Result {
	user := c.getUser(email)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = email
			c.Session["user_id"] = strconv.Itoa(user.ID)
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + email)

			// Set language from user's preference.
			// TODO add language field to User model
			// language field should be string and ISO639-1 codes.
			// Region field should be string and ISO3166-1 alpha-2 code
			
			cookie := http.Cookie{Name: "REVEL_LANG", Value: "en-US", Path: "/"}
			c.SetCookie(&cookie)
			
			user.LastLoginAt = time.Now()
			c.Tx.Save(&user)
			return c.Redirect(routes.Hello.Welcome())
		}
	}

	c.Flash.Out["user"] = email
	c.Flash.Error("Login failed")
	return c.Redirect(routes.Application.Index())
}

func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Index())
}
*/
