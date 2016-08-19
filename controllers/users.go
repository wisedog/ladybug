package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

// UserProfile renders first page of user's profile
func UserProfile(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	var targetUser models.User
	if err := c.Db.Where("id = ?", id).First(&targetUser); err.Error != nil {
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "An error is occurred while get target user"}
	}

	var activities []models.Activity
	if err := c.Db.Where("user_id = ?", id).Preload("User").Find(&activities); err.Error != nil {
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error is occurred while get activities."}
	}

	idNum, _ := strconv.Atoi(id)
	isSameID := 0
	if idNum == c.User.ID {
		isSameID = 1
	}

	items := map[string]interface{}{
		"Activities":   activities,
		"Active_idx":   isSameID,
		"TargetUser":   targetUser,
		"ShowUserMenu": true,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/users/profile.tmpl")
}

// UserGetNameList returns users' ID and name.
func UserGetNameList(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	var users []models.User
	if err := c.Db.Find(&users); err.Error != nil {
		log.Error("User", "msg", err.Error)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "An error is occurred while get users"}
	}

	type data struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Photo string `json:"photo"`
	}

	renderData := make([]data, len(users))

	for i, x := range users {
		renderData[i].ID = x.ID
		renderData[i].Name = x.Name
		renderData[i].Email = x.Email
		renderData[i].Photo = x.Photo

	}
	return RenderJSON(w, renderData)
}

// UserUpdateProfile is a POST handler for updating user's profile
func UserUpdateProfile(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	idInt, _ := strconv.Atoi(id)

	var targetUser models.User
	if err := c.Db.Where("id = ?", id).First(&targetUser); err.Error != nil {
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Bad Request"}
	}

	if c.User.ID != idInt {
		// if user does not have the authorization to modify data, return unauth
		if c.User.Role != models.RoleAdmin {
			return errors.HttpError{Status: http.StatusUnauthorized, Desc: "Unauthorized update user's profile"}
		}
	}

	// parse data from POST form
	var user models.User

	if err := r.ParseForm(); err != nil {
		log.Error("Users", "type", "http", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "ParseForm failed"}
	}

	if err := c.Decoder.Decode(&user, r.PostForm); err != nil {
		log.Warn("Build", "Error", err, "msg", "Decode failed but go ahead")
	}

	// only 4 fields - name, email, location, notes could be handled
	targetUser.Name = user.Name
	//NOT_NOW targetUser.Email = user.Email
	targetUser.Location = user.Location
	targetUser.Notes = user.Notes

	// validate it!

	log.Debug("Users", "After decode", user)

	// update
	if err := c.Db.Save(&targetUser); err.Error != nil {
		log.Debug("Users", "type", "db", "msg", "update user operation is failed")
		return nil
	}

	// TODO if email is changed, session information should be changed!
	//return RenderJSONWithStatus(w, Resp{Msg: "OK"}, http.StatusOK)
	http.Redirect(w, r, "/user/profile/"+id, http.StatusFound)
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
