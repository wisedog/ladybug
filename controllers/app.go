package controllers

import (
	"fmt"
	"os"
	"strings"

	"encoding/gob"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/robfig/config"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"
	"golang.org/x/crypto/bcrypt"

	log "gopkg.in/inconshreveable/log15.v2"
)

// Costant message
const (
	LadybugSession   = "ladybug_session"
	FlashAuthMessage = "flash_auth_msg"
)

// function map used in template files
var (
	funcMap = template.FuncMap{
		"nl2br": func(str string) string { return strings.Replace(str, "\n", "<br />", -1) },
		"msg": func(key string) string {
			if messages == nil {
				messages = LoadI18nMessage()
			}
			s, _ := messages.String("", key)
			return s
		},
		"categoryi18n": func(key string) string {
			return ""
		},
	}
	messages *config.Config
)

func init() {

	// for flash messages
	gob.Register(&models.Build{})
	gob.Register(&models.Project{})
	gob.Register(&models.Requirement{})
	gob.Register(&map[string]string{})

	// TODO store category map to static variable.
	// in funcMap - categorui18n function refers the var.
}

// connected is private utilitiy function for checking
// this user id now on connected or not.
func connected(c *interfacer.AppContext, r *http.Request) *models.User {
	session, err := c.Store.Get(r, "ladybug")
	if err != nil {
		log.Warn("An error in connected(),", "msg", err)
		return nil
	}
	email := session.Values["user"]
	uid := session.Values["uid"]

	var u models.User
	c.Db.Where("email = ? and id = ?", email, uid).First(&u)

	if u.ID == 0 {
		return nil
	}
	return &u
}

func renderTemplate(w http.ResponseWriter, tmpl string, msg string, isLogout bool) {
	dir, _ := os.Getwd()
	target := filepath.Join(dir, "views", "application", tmpl)
	t, err := template.ParseFiles(target)

	if err != nil {
		log.Error("App", "msg", "Template parse files error", "err", err)
	}

	items := map[string]interface{}{
		"Messages": msg,
		"IsLogout": isLogout,
	}

	fmt.Println("items", items)

	if err := t.Execute(w, items); err != nil {
		log.Error("App", "msg", "Template Execute error", "err", err)
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
func LoginPage(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
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
		c.Db.Where("email = ? and id = ?", email, uid).First(&u)

		if u.ID == 0 {
			return nil //TODO return unauth
		}

		http.Redirect(w, r, "/hello", http.StatusFound)

	} else {
		log.Debug("App", "msg", "not login, rendering login page")

		// Get the previously flashes, if any.
		var s string
		var logoutFlag bool

		if msg := session.Flashes(FlashAuthMessage); len(msg) > 0 {
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
func LogOut(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	session, err := c.Store.Get(r, "ladybug")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return err
	}
	session.Values["user"] = ""
	session.Values["uid"] = 0
	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
	log.Debug("App", "msg", "logout")
	session.AddFlash("You just logged out.", FlashAuthMessage)

	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

// Register renders just a user registration page
func Register(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	dir, _ := os.Getwd()
	target := filepath.Join(dir, "views", "application", "register.html")
	t, err := template.ParseFiles(target)

	if err != nil {
		log.Error("App", "msg", "Template parse files error", "err", err)
		return err
	}

	session, e := c.Store.Get(r, "ladybug")
	var errorMap map[string]string
	if e == nil {
		errorMap = getErrorMap(session)
		if errorMap == nil {
			log.Info("app", "msg", "map is nil")
			errorMap = make(map[string]string)
		}
	} else {
		errorMap = make(map[string]string)
	}

	items := map[string]interface{}{
		"Messages": "msg",
		"ErrorMap": errorMap,
	}

	if err := t.Execute(w, items); err != nil {
		log.Error("App", "msg", "Template Execute error", "err", err)
		return err
	}

	return nil
}

// SaveUser is a POST handler from Register page
func SaveUser(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	// First parse form value
	if err := r.ParseForm(); err != nil {
		log.Error("User", "type", "http", "msg ", err)
	}

	password1 := r.FormValue("Password1")
	password2 := r.FormValue("Password2")

	// check those two password is identical
	if password1 != password2 {
		session, _ := c.Store.Get(r, "ladybug")
		errorMap := make(map[string]string)
		errorMap["Password"] = "The passwords are not identical."
		session.AddFlash(errorMap, ErrorMsg)

		session.Save(r, w)

		http.Redirect(w, r, "/register", http.StatusFound)
		return nil
	}

	email := r.FormValue("Email")
	name := r.FormValue("Name")

	log.Debug("debug", "email", email, "name", name)

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(password1), bcrypt.DefaultCost)

	usr := models.User{Email: email, Name: name, Password: password1, HashedPassword: hashedPassword}
	log.Debug("debug", "hashed", hashedPassword)

	// Validate input value
	// if validate failed, redirect register page with Flash
	if errorMap := usr.Validate(); len(errorMap) > 0 {
		log.Debug("debug", "msg", "in validation fail", "errormap", errorMap)

		session, e := c.Store.Get(r, "ladybug")
		if e != nil {
			log.Warn("error ", "msg", e.Error)
		}

		session.AddFlash(usr, "LADYBUG_USER")
		session.AddFlash(errorMap, ErrorMsg)

		session.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return nil
	}

	var tmpUsr models.User
	// check duplicated email
	if err := c.Db.Where("email = ?", email).First(&tmpUsr).Error; err == nil {
		log.Debug("debug", "msg", "Found duplicated")
		// Found duplicated
		errorMap := make(map[string]string)
		session, e := c.Store.Get(r, "ladybug")
		if e != nil {
			log.Warn("error ", "msg", e.Error)
		}

		errorMap["Email"] = "Duplicated Email"

		session.AddFlash(usr, "LADYBUG_USER")
		session.AddFlash(errorMap, ErrorMsg)

		session.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return nil
	}

	// save to db
	c.Db.NewRecord(usr)

	if err := c.Db.Create(&usr).Error; err != nil {
		log.Error("App", "Type", "database", "Error", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "could not create a build model"}
	}

	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}

// LogAndHTTPError makes log messages and return http error with given status http code
// default log level is ERROR
func LogAndHTTPError(status int, module string, errType string, msg string) error {
	log.Error(module, "type", errType, "msg", msg)
	return errors.HttpError{Status: status, Desc: msg}
}
