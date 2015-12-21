package controllers

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"github.com/wisedog/ladybug/app/routes"
	/*"strconv"
	"golang.org/x/crypto/bcrypt"
	"strings"*/)

type Hello struct {
	Application
}

func (c Hello) Index() revel.Result {
	var user *models.User

	if user = c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}

	/*TODO Read all suites and testcases of the project
	    and response via json model.
		JSTree component will receive the response and build the tree.
		To recude response size, return only ID, summary, Automatic/Manual info
	*/
	var prjs []models.Project
	c.Tx.Find(&prjs)

	revel.INFO.Println("AAA:", prjs)

	return c.Render(prjs, user)

}
