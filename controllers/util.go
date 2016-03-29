package controllers

import (
  "net/http"
  "encoding/json"
  "html/template"

  "github.com/gorilla/sessions"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/errors"
  "github.com/wisedog/ladybug/interfacer"
  log "gopkg.in/inconshreveable/log15.v2"
)

// getPriorityI18n function returns localization string 
func getPriorityI18n(priority int) string{
	str := "Unknown Status"
	switch priority{
		case models.PriorityHighest:
			str = GetI18nMessage("priority.highest")
		case models.PriorityHigh:
      str = GetI18nMessage("priority.high")
		case models.PriorityMedium:
			str = GetI18nMessage("priority.medium")
		case models.PriorityLow:
			str = GetI18nMessage("priority.low")
		case models.PriorityLowest:
			str = GetI18nMessage("priority.lowest")
	}
	
	return str
}

func getErrorMap(session *sessions.Session) *map[string]string{
  if fm := session.Flashes(ErrorMsg); fm != nil {
    b, ok := fm[0].(*map[string]string)
    if ok{
      delete(session.Values, ErrorMsg)
      return b
    }else{
      log.Debug("Build", "msg", "flash type assertion failed")
    }    
  }
  return nil
}

// Render2 renders templates with structure-typed interface data
func Render(w http.ResponseWriter, data interface{},  templates ...string) error{
  t, err := template.New("base.tmpl").Funcs(funcMap).ParseFiles(templates...)

  if err != nil{
    log.Error("Util", "type", "rendering", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "Template ParseFiles error"}
  }

  if err = t.Execute(w, data); err != nil{
    log.Error("Util", "type", "rendering", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "Template Exection Error"}
  }

  return nil
}

// Render2 renders templates with map-typed interface data
func Render2(c *interfacer.AppContext, w http.ResponseWriter, data interface{},  templates ...string) error{
  t, err := template.New("base.tmpl").Funcs(funcMap).ParseFiles(templates...)

  if err != nil{
    log.Error("Util", "type", "rendering", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "Template ParseFiles error"}
  }

  item := data.(map[string]interface{})
  item["User"] = c.User
  item["ProjectName"] = c.ProjectName

  if err = t.Execute(w, item); err != nil{
    log.Error("Util", "type", "rendering", "msg ", err )
    return errors.HttpError{http.StatusInternalServerError, "Template Exection Error"}
  }

  return nil
}

// RenderJson renders JSON 
func RenderJson(w http.ResponseWriter, data interface{}) error{
  js, err := json.Marshal(data)
  if err != nil {
    log.Error("Builds", "msg", "Json Marshalling failed in ValidateTool")
    return err
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
  return nil
}