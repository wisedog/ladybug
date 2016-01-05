package controllers

import (
	"github.com/revel/revel"
	)


type Specifications struct {
	Application
}

func (c Specifications) Spec(project string) revel.Result {
	return c.Render(project)
}