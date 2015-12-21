package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
)

type Users struct {
	Application
}

func (c Users) Index(id int) revel.Result {
	var u models.User

	revel.INFO.Println("In Cookie:", c.Session["user"])
	c.Tx.Where("id = ?", id).First(&u)

	return c.Render(u)

}

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
	var usr models.User
	r := c.Tx.Where("id = ?", id).First(&usr)

	if r.Error != nil {
		return c.NotFound("Something wrong....")
	}

	return c.Render(usr)
}
