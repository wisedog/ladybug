package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"
	log "gopkg.in/inconshreveable/log15.v2"
)

// getPriorityI18n function returns localization string
func getPriorityI18n(priority int) string {
	str := "Unknown Status"
	switch priority {
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

func getReqStatusI18n(reqStatus int) string {
	str := "Unknown Status"
	switch reqStatus {
	case models.ReqStatusRework:
		str = GetI18nMessage("requirement.status.rework")
	case models.ReqStatusDraft:
		str = GetI18nMessage("requirement.status.draft")
	case models.ReqStatusFinished:
		str = GetI18nMessage("requirement.status.finished")
	case models.ReqStatusInReview:
		str = GetI18nMessage("requirement.status.inreview")
	case models.ReqStatusNotTestable:
		str = GetI18nMessage("requirement.status.nottestable")
	case models.ReqStatusDeprecated:
		str = GetI18nMessage("requirement.status.deprecated")
	}

	return str
}

func getErrorMap(session *sessions.Session) map[string]string {
	if fm := session.Flashes(ErrorMsg); fm != nil {
		b, ok := fm[0].(*map[string]string)
		if ok {
			delete(session.Values, ErrorMsg)
			return *b
		} else {
			log.Debug("Build", "msg", "flash type assertion failed")
		}
	}
	return nil
}

// Render renders templates with structure-typed interface data
func Render(w http.ResponseWriter, data interface{}, templates ...string) error {
	t, err := template.New("base.tmpl").Funcs(funcMap).ParseFiles(templates...)

	if err != nil {
		log.Error("Util", "type", "rendering", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Template ParseFiles error"}
	}

	if err = t.Execute(w, data); err != nil {
		log.Error("Util", "type", "rendering", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Template Exection Error"}
	}

	return nil
}

// Render2 renders templates with map-typed interface data
func Render2(c *interfacer.AppContext, w http.ResponseWriter, data interface{}, templates ...string) error {
	t, err := template.New("base.tmpl").Funcs(funcMap).ParseFiles(templates...)

	if err != nil {
		log.Error("Util", "type", "rendering", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Template ParseFiles error"}
	}

	item := data.(map[string]interface{})
	item["User"] = c.User
	item["ProjectName"] = c.ProjectName

	if err = t.Execute(w, item); err != nil {
		log.Error("Util", "type", "rendering", "msg ", err)
		return errors.HttpError{Status: http.StatusInternalServerError, Desc: "Template Exection Error"}
	}

	return nil
}

// RenderJSONWithStatus renders JSON format with data with status specified
func RenderJSONWithStatus(w http.ResponseWriter, data interface{}, statusCode int) error {
	js, err := json.Marshal(data)
	if err != nil {
		log.Error("Builds", "msg", "Json Marshalling failed in ValidateTool")
		return err
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

// RenderJSON converts data interface to JSON and renders JSON format
func RenderJSON(w http.ResponseWriter, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		log.Error("Builds", "msg", "Json Marshalling failed in ValidateTool")
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

// ValidationFailAndRedirect store errors and redirect if there are error messages
func ValidationFailAndRedirect(c *interfacer.AppContext, w http.ResponseWriter,
	r *http.Request, errorMap map[string]string, url string, value interface{}) {
	if len(errorMap) > 0 {
		log.Debug("Validation failed")

		session, e := c.Store.Get(r, "ladybug")
		if e != nil {
			log.Warn("error ", "msg", e.Error)
		}

		session.AddFlash(value, "LADYBUG_BUILD")
		session.AddFlash(errorMap, ErrorMsg)

		session.Save(r, w)
		http.Redirect(w, r, url, http.StatusFound)
	}

}
