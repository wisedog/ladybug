package controllers

import (
	"net/http"

	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

// SettingGeneral renders general setting of project
// target : description, name(Not ID), Access level .....
func SettingGeneral(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {

	// get project ID
	var prj models.Project
	if err := c.Db.Where("name = ?", c.ProjectName).First(&prj).Error; err != nil {
		log.Error("Setting", "msg", err.Error(), "additional", "not found project", "target", c.ProjectName)
		return errors.HttpError{Status: http.StatusBadRequest, Desc: "Not found project"}
	}

	// Only manager and admin can access this menu. check auth level from role table
	var role models.Role
	if err := c.Db.Where("project_id = ?", prj.ID).Where("user_id = ?", c.User.ID).First(&role).Error; err != nil {
		log.Error("Setting", "msg", err.Error(), "project_id", prj.ID, "user_id", c.User.ID)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Not found role in role table"}
	}

	if role.UserRole > models.RoleManager {
		log.Error("Setting", "msg", "Not enough permission to access", "additional", role.UserRole)
		return errors.HttpError{Status: http.StatusUnauthorized, Desc: "Not enough permission to access"}
	}

	items := map[string]interface{}{
		"Active_idx":         9,
		"ShowProjectSetting": true,
		"Project":            prj,
	}

	return Render2(c, w, items, "views/base.tmpl", "views/projectsetting/settinggeneral.tmpl")
}
