package controllers

import (
  "fmt"
  "os"

  "net/http"
  "html/template"
  "path/filepath"

  "golang.org/x/crypto/bcrypt"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  //"github.com/wisedog/ladybug/errors"
  "github.com/gorilla/securecookie"
  "github.com/gorilla/sessions"
  //"github.com/jinzhu/gorm"	
   log "gopkg.in/inconshreveable/log15.v2"
)

const (
  LADYBUG_SESSION = "ladybug_session"
)

funcMap := template.FuncMap {
  "nl2br": func (str string) string { return strings.Replace(str, "\n", "<br />", -1) },
  "msg" : func (str string) string{
      // TODO. much to do
      return str
    },
  }

//@deprecated
var store = sessions.NewCookieStore([]byte("something-very-secret"))

var cookieHandler = securecookie.New(
  securecookie.GenerateRandomKey(64),
  securecookie.GenerateRandomKey(32),
  )

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
    return nil  //TODO return unauth
  }
  return &u
}

//@deprecated
func getUserId(request *http.Request) (userName string) {
  if cookie, err := request.Cookie(LADYBUG_SESSION); err == nil {
    cookieValue := make(map[string]string)
    if err = cookieHandler.Decode(LADYBUG_SESSION, cookie.Value, &cookieValue); err == nil {
      userName = cookieValue["name"]
    }
  }
  return userName
}

func clearSession(response http.ResponseWriter){
  cookie := &http.Cookie{
    Name:   LADYBUG_SESSION,
    Value:  "",
    Path:   "/",
    MaxAge: -1,
  }
  http.SetCookie(response, cookie)
}

//@deprecated
func setSession(userId string, response http.ResponseWriter){
  value := map[string]string{
    "name": userId,
  }
  if encoded, err := cookieHandler.Encode("session", value); err == nil {
    cookie := &http.Cookie{
      Name:  "session",
      Value: encoded,
      Path:  "/",
    }
    http.SetCookie(response, cookie)
  }
}

func renderTemplate(w http.ResponseWriter, tmpl string, Msg string) {
  dir, _ := os.Getwd()
  target := filepath.Join(dir, "views", "application", tmpl )
  fmt.Println(target)
  t, _ := template.ParseFiles(target)
  
  t.Execute(w, Msg)
}

//@deprecated
func getUser(c *interfacer.AppContext, email string) *models.User{
  var u models.User
  err := c.Db.Where("email = ?", email).First(&u)
  if err.Error != nil{
      // TODO
  }
  return &u
}

// Login checks user email, password 
func Login(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

  fmt.Println("login Processing")
  session, err := c.Store.Get(r, "ladybug")
  if err != nil {

      http.Error(w, err.Error(), 500)
      return nil
  }

  email := r.FormValue("email")
  passwd := r.FormValue("password")
  user := getUser(c, email)
  if user != nil {
    err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(passwd))
    if err == nil {
      session.Values["user"] = email
      session.Values["uid"] = user.ID

      fmt.Println("user", user)
      fmt.Println("uid", user.ID)

      // Save it before we write to the response/return from the handler.
      session.Save(r, w)

      http.Redirect(w, r, "/hello", http.StatusFound)
    }
  }
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
    fmt.Println("Already logged in")
    fmt.Println("in Loginpage user", session.Values["user"])
    fmt.Println("in Loginpage uid ", session.Values["uid"])
    email := session.Values["user"]
    uid := session.Values["uid"]

    var u models.User
    c.Db.Where("email = ? and id = ?", email, uid ).First(&u)

    if u.ID == 0{
      return nil  //TODO return unauth
    }

    http.Redirect(w, r, "/hello", http.StatusFound)

  }else{
    fmt.Println("not login, rendering login page")
    // Get the previously flashes, if any.

    var s string
    /* gorilla/session's flash is not working now FIXME
    fmt.Println("flashes : ", session.Flashes("msg"))
    
    if msg := session.Flashes(); len(msg) > 0{
      fmt.Println("Flash : ", msg)
      v, ok := msg[0].(string)
      if ok {
        s = v
      }
    }*/
    
    renderTemplate(w, "index.html", s)
    return nil
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

  // Set a new flash.... is not working.. FIXME
  // session.AddFlash("You just logged out.", "msg")

  session.Save(r, w)
  http.Redirect(w, r, "/login", http.StatusFound)
  return nil
}

/*
func (c App) LoginPage(w http.ResponseWriter, r *http.Request) {
  fmt.Println("login page")
  u := &models.User{}
  c.renderTemplate(w, "index.html", u)

}


func (c App) Logout(w http.ResponseWriter, r *http.Request) {
  clearSession(w)
  http.Redirect(w, r, "/", 302)
}

/*
func (c Application) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Hello.Welcome())
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c Application) Register(){
	return nil
}
*/
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
