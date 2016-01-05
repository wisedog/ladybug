package controllers

import (
	"github.com/revel/revel"
	)


type Reports struct {
	Application
}

func (c Reports) Reports(project string) revel.Result {
	return c.Render(project)
}