package controllers

import (
  "net/http"

  //"github.com/gorilla/mux"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"

  log "gopkg.in/inconshreveable/log15.v2"
)

// UserProfile renders first page of user's profile
func UserProfile(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  log.Debug("User", "msg", "in UserProfile")

  var activities []models.Activity
  if err := c.Db.Where("user_id = ?", c.User.ID).Preload("User").Find(&activities); err.Error != nil{
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "An error is occurred while get activities."}
  }

  items := map[string]interface{}{
    "Activities" : activities,
    "Active_idx" : 0,
  }


  return Render2(c, w, items, "views/base.tmpl", "views/users/profile.tmpl")
}

// UserGetNameList returns users' ID and name.
func UserGetNameList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  var users []models.User
  if err := c.Db.Find(&users); err.Error != nil{
    log.Error("User", "msg", err.Error)
    return errors.HttpError{Status : http.StatusInternalServerError, Desc : "An error is occurred while get users"}
  }

  type data struct{
    ID  int `json:"id"`
    Name  string  `json:"name"`
    Email string  `json:"email"`
    Photo string  `json:"photo"`
  }

  renderData := make([]data, len(users))

  for i, x := range users{
    renderData[i].ID = x.ID
    renderData[i].Name = x.Name
    renderData[i].Email = x.Email
    renderData[i].Photo = x.Photo

  }
  return RenderJSON(w, renderData)
}

// UserUpdateProfile is a POST handler for updating user's profile
func UserUpdateProfile(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

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

*/