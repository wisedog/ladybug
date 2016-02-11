package buildtools

import (
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
)

/*
type IController interface {
  PrintMessage()
}

type BuildTool struct {
  parent IController
}

func (bt *BuildTool) PrintParentMessage() {
  bt.parent.PrintMessage()
}

func NewChild(parent IController) *BuildTool {
  return &BuildTool{parent: parent }
}
*/

type Travis struct {
}

func (t Travis) ConnectTest(id int) int{
	var issues []models.Issue

	revel.INFO.Println("In Travis", id, issues)
	return id +10
}

/*
func Detail(project string, id int, ) revel.Result {
	var issues []models.Issue

	revel.INFO.Println("AAA", project, id)
	return c.Render(issues)
}
*/