package controllers

import (
  "net/http"
	
  "github.com/gorilla/mux"
	"github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
  log "gopkg.in/inconshreveable/log15.v2"
	)

func Reports(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
  var user *models.User
  log.Debug("Reports", "msg", "in Reports")
  if user = connected(c, r); user == nil {
    log.Info("Not found login information.")
    http.Redirect(w, r, "/", http.StatusFound)

    return nil
  }

  vars := mux.Vars(r)
  projectName := vars["projectName"]

  items := map[string]interface{}{
    "User": user,
    "ProjectName" : projectName,
    "Active_idx" : 7,
  }

	return Render(w, items, "views/base.tmpl", "views/testdesign/designindex.tmpl")
}