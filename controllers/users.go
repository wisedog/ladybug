package controllers

import (
	"fmt"
  "net/http"
  "html/template"

  "github.com/gorilla/mux"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"

  log "gopkg.in/inconshreveable/log15.v2"
)

func UserProfile(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  log.Debug("User", "msg", "in UserProfile")
  vars := mux.Vars(r)
  id := vars["id"]

  var user models.User

  err := c.Db.Where("id = ?", id).First(&user)
  if err.Error != nil{
    return errors.HttpError{http.StatusBadRequest, "Bad request"}
  }

  var activities []models.Activity
  if err := c.Db.Where("user_id = ?", user.ID).Preload("User").Find(&activities); err.Error != nil{
    return errors.HttpError{http.StatusInternalServerError, "An error is occurred while get activities."}
  }

  items := struct {
    User *models.User
    Activities  []models.Activity  
    Active_idx  int
  }{
    User: &user,
    Activities : activities,
    Active_idx : 0,
  }

  t, er := template.New("base.tmpl").Funcs(funcMap).ParseFiles(
    "views/base.tmpl",
    "views/users/profile.tmpl",
    )


  if er != nil{
    fmt.Println("Err ", er )
    return er
  }
  
  if err := t.Execute(w, items); err != nil{
    fmt.Println("Execution failed : ", err)
    return err
  }
  
  return nil
}

func UserGeneral(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  
  return nil
}
/*
func (c Users) SaveGeneral(id int, Name string, Language string) revel.Result {

	usr := *(c.connected())
	if usr.ID != id {
		return c.NotFound("Something wrong....")
	}

	usr.Name = Name
	usr.Language = Language
	usr.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Users.Index(id))

	}

	c.Flash.Data["Language"] = usr.Language
	c.Tx.Save(&usr)

	return c.Redirect(routes.Users.Index(id))

}

func (c Users) Register() revel.Result {

	return c.Render()
}

func (c Users) Edit(id int) revel.Result {

	return c.Render()
}

func (c Users) Profile(id int) revel.Result {
	var user models.User
	r := c.Tx.Where("id = ?", id).First(&user)

	if r.Error != nil {
		return c.NotFound("Something wrong....")
	}
	
	var activities []models.Activity
	c.Tx.Where("user_id = ?", user.ID).Preload("User").Find(&activities)

	return c.Render(user, activities)
}
*/