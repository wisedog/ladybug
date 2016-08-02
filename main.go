package main

import (
	"flag"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/wisedog/ladybug/controllers"
	"github.com/wisedog/ladybug/database"
	"github.com/wisedog/ladybug/errors"
	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

const (
	CookieKey = "ladybug"
)

// connected2 is private utilitiy function for checking
// this user id now on connected or not.
func connected2(c *interfacer.AppContext, r *http.Request) *models.User {
	session, err := c.Store.Get(r, "ladybug")
	if err != nil {
		log.Warn("An error in connected(),", "msg", err)
		return nil
	}
	email := session.Values["user"]
	uid := session.Values["uid"]

	var u models.User
	c.Db.Where("email = ? and id = ?", email, uid).First(&u)

	if u.ID == 0 {
		return nil
	}
	return &u
}

type handlerFunc func(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error

//
func myHandler(c *interfacer.AppContext, h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errors.HandleError(w, r, func() error {
			return h(c, w, r)
		}())
	}
}

// authHandler checks authorization from session
func authHandler(c *interfacer.AppContext, h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errors.HandleError(w, r, func() error {
			var user *models.User
			if user = connected2(c, r); user == nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return nil
			}
			c.User = user
			vars := mux.Vars(r)
			c.ProjectName = vars["projectName"]

			return h(c, w, r)
		}())
	}
}

// notFound handles 404 error and render custom page.
// see views/errors/404.tmpl
func notFound(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error {
	t, err := template.New("base.tmpl").ParseFiles("views/base.tmpl", "views/errors/404.tmpl")

	if err != nil {
		log.Error("Util", "type", "rendering", "msg ", err)
	}
	var user *models.User
	if user = connected2(c, r); user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	item := map[string]interface{}{
		"Active_idx": 0,
		"User":       c.User,
	}

	if err = t.Execute(w, item); err != nil {
		log.Error("Util", "type", "rendering", "msg ", err)
	}

	return nil
}

// main is entry point of this application
func main() {

	log.Info("App initialize...")

	// create router with gorilla/mux
	r := mux.NewRouter()

	// create application context
	ctx := &interfacer.AppContext{}

	// parse program argument
	modePtr := flag.String("mode", "", "mode")
	portPtr := flag.Int("port", 8000, "(optional)binding port. default value is 8000")
	addrPtr := flag.String("addr", "localhost", "(optional)binding address. default value is localhost")
	databasePtr := flag.String("db", "", "(optional)database url(dialect)")
	flag.Parse()
	log.Info("APP", "Starting Mode", *modePtr)

	// load config
	ctx.Config = interfacer.LoadConfigWithArgs(*modePtr, *addrPtr, *portPtr, *databasePtr, "./ladybug.conf")

	log.Info("Initialize Database...")
	if db, err := database.InitDB(ctx.Config); err != nil {
		log.Crit("Database initialization is failed.")
		return
	} else {
		ctx.Db = db
	}

	ctx.Store = sessions.NewCookieStore([]byte("ladybug"))
	ctx.Decoder = schema.NewDecoder()

	// load i18n resources
	controllers.LoadI18nMessage()

	//// Routes consist of a path and a handler function.

	// application overall
	r.HandleFunc("/", myHandler(ctx, controllers.LoginPage)).Methods("GET")
	r.HandleFunc("/login", myHandler(ctx, controllers.LoginPage)).Methods("GET")
	r.HandleFunc("/login", myHandler(ctx, controllers.Login)).Methods("POST")
	r.HandleFunc("/logout", myHandler(ctx, controllers.LogOut)).Methods("GET")
	r.HandleFunc("/register", myHandler(ctx, controllers.Register)).Methods("GET")
	r.HandleFunc("/hello", authHandler(ctx, controllers.Welcome)).Methods("GET")
	r.HandleFunc("/saveuser", myHandler(ctx, controllers.SaveUser)).Methods("POST")

	manage := r.PathPrefix("/manage/").Subrouter()

	// for project creation
	manage.HandleFunc("/project/create", authHandler(ctx, controllers.ProjectCreate)).Methods("GET")
	manage.HandleFunc("/project/save", authHandler(ctx, controllers.ProjectSave)).Methods("POST")

	// define user subrouter
	user := r.PathPrefix("/user/").Subrouter()
	user.HandleFunc("/profile/{id:[0-9]+}", authHandler(ctx, controllers.UserProfile)).Methods("GET")
	user.HandleFunc("/profile/update/{id:[0-9]+}", authHandler(ctx, controllers.UserUpdateProfile)).Methods("POST")
	//user.HandleFunc("/general/{id:[0-9]+}", authHandler(ctx, controllers.UserGeneral)).Methods("GET")
	//user.HandleFunc("/get/list", authHandler(ctx, controllers.UserGetNameList)).Methods("GET")
	user.HandleFunc("/get/list", authHandler(ctx, controllers.UserGetNameList)).Methods("POST")

	// define project subrouter
	project := r.PathPrefix("/project/").Subrouter()

	// TODO add reserved word create, save, list, get. those words are not allowed to be project name
	project.HandleFunc("/get/list", authHandler(ctx, controllers.GetProjectList)).Methods("GET")

	// project specific
	project.HandleFunc("/{projectName}", authHandler(ctx, controllers.Dashboard)).Methods("GET")
	project.HandleFunc("/{projectName}/dashboard", authHandler(ctx, controllers.Dashboard)).Methods("GET")
	project.HandleFunc("/{projectName}/design", authHandler(ctx, controllers.DesignIndex)).Methods("GET")
	project.HandleFunc("/{projectName}/build", authHandler(ctx, controllers.BuildsIndex)).Methods("GET")
	project.HandleFunc("/{projectName}/req", authHandler(ctx, controllers.RequirementIndex)).Methods("GET")
	project.HandleFunc("/{projectName}/testplan", authHandler(ctx, controllers.PlanIndex)).Methods("GET")
	project.HandleFunc("/{projectName}/exec", authHandler(ctx, controllers.ExecIndex)).Methods("GET")
	project.HandleFunc("/{projectName}/milestone", authHandler(ctx, controllers.MilestoneIndex)).Methods("GET")

	// section
	section := project.PathPrefix("/{projectName}/section/").Subrouter()
	section.HandleFunc("/testcase/{sectionID}", authHandler(ctx, controllers.GetAllTestCases)).Methods("GET")
	section.HandleFunc("/add/{sectionID}", authHandler(ctx, controllers.SectionAdd)).Methods("GET")

	// test cases
	testcase := project.PathPrefix("/{projectName}/case/").Subrouter()
	testcase.HandleFunc("/view/{id:[0-9]+}", authHandler(ctx, controllers.CaseView)).Methods("GET")
	testcase.HandleFunc("/edit/{id:[0-9]+}", authHandler(ctx, controllers.CaseEdit)).Methods("GET")
	testcase.HandleFunc("/update/{id:[0-9]+}", authHandler(ctx, controllers.CaseUpdate)).Methods("POST")
	testcase.HandleFunc("/add", authHandler(ctx, controllers.CaseAdd)).Methods("GET")
	testcase.HandleFunc("/save", authHandler(ctx, controllers.CaseSave)).Methods("POST")
	testcase.HandleFunc("/delete", authHandler(ctx, controllers.CaseDelete)).Methods("POST")

	// builds
	build := project.PathPrefix("/{projectName}/build").Subrouter()
	build.HandleFunc("/", authHandler(ctx, controllers.BuildsIndex)).Methods("GET")
	build.HandleFunc("/add", authHandler(ctx, controllers.BuildsAddProject)).Methods("GET")
	build.HandleFunc("/edit/{id:[0-9]+}", authHandler(ctx, controllers.BuildsEditProject)).Methods("GET")
	build.HandleFunc("/view/{id:[0-9]+}", authHandler(ctx, controllers.BuildsView)).Methods("GET")
	build.HandleFunc("/save", authHandler(ctx, controllers.BuildsSaveProject)).Methods("POST")
	build.HandleFunc("/integrate", authHandler(ctx, controllers.Integrate)).Methods("GET")
	build.HandleFunc("/tool/validate", authHandler(ctx, controllers.ValidateTool)).Methods("POST")
	build.HandleFunc("/tool/add", authHandler(ctx, controllers.AddTool)).Methods("POST")

	// TODO Build - Update

	build.HandleFunc("/list/{id:[0-9]+}", authHandler(ctx, controllers.GetBuildItems)).Methods("GET")
	build.HandleFunc("/additem/{id:[0-9]+}", authHandler(ctx, controllers.BuildsAddItem)).Methods("GET")
	build.HandleFunc("/saveitem", authHandler(ctx, controllers.BuildsSaveItem)).Methods("POST")
	//TODO viewitem

	// test plan
	plan := project.PathPrefix("/{projectName}/testplan").Subrouter()
	plan.HandleFunc("/", authHandler(ctx, controllers.PlanIndex)).Methods("GET")
	plan.HandleFunc("/add", authHandler(ctx, controllers.PlanAdd)).Methods("GET")
	plan.HandleFunc("/save", authHandler(ctx, controllers.PlanSave)).Methods("POST")
	plan.HandleFunc("/delete", authHandler(ctx, controllers.PlanDelete)).Methods("POST")
	plan.HandleFunc("/view/{id:[0-9]+}", authHandler(ctx, controllers.PlanView)).Methods("GET")
	plan.HandleFunc("/run/{id:[0-9]+}", authHandler(ctx, controllers.PlanRun)).Methods("GET")

	// requirements
	req := project.PathPrefix("/{projectName}/req").Subrouter()
	req.HandleFunc("/", authHandler(ctx, controllers.RequirementIndex)).Methods("GET")
	req.HandleFunc("/add", authHandler(ctx, controllers.AddRequirement)).Methods("GET")
	req.HandleFunc("/view/{id:[0-9]+}", authHandler(ctx, controllers.ViewRequirement)).Methods("GET")
	req.HandleFunc("/edit/{id:[0-9]+}", authHandler(ctx, controllers.EditRequirement)).Methods("GET")
	req.HandleFunc("/list/{id:[0-9]+}", authHandler(ctx, controllers.RequirementList)).Methods("GET")
	req.HandleFunc("/save/{id:[0-9]+}", authHandler(ctx, controllers.SaveRequirement)).Methods("POST")
	req.HandleFunc("/delete/{id:[0-9]+}", authHandler(ctx, controllers.DeleteRequirement)).Methods("POST")

	// testexec
	exec := project.PathPrefix("/{projectName}/exec").Subrouter()
	exec.HandleFunc("/", authHandler(ctx, controllers.ExecIndex)).Methods("GET")
	exec.HandleFunc("/run/{id:[0-9]+}", authHandler(ctx, controllers.ExecRun)).Methods("GET")
	exec.HandleFunc("/remove", authHandler(ctx, controllers.ExecRemove)).Methods("POST")
	exec.HandleFunc("/deny", authHandler(ctx, controllers.ExecDeny)).Methods("POST")
	exec.HandleFunc("/update", authHandler(ctx, controllers.ExecUpdateResult)).Methods("POST")
	exec.HandleFunc("/done", authHandler(ctx, controllers.ExecDone)).Methods("POST")

	// milestone
	milestone := project.PathPrefix("/{projectName}/milestone").Subrouter()
	milestone.HandleFunc("/", authHandler(ctx, controllers.MilestoneIndex)).Methods("GET")
	milestone.HandleFunc("/add", authHandler(ctx, controllers.MilestoneAdd)).Methods("GET")

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/",
		http.FileServer(http.Dir("public/"))))
	r.NotFoundHandler = http.HandlerFunc(authHandler(ctx, notFound))

	log.Info("APP", "Binding Address", ctx.Config.GetBindAddress())
	// Bind to a port and pass our router in
	http.ListenAndServe(ctx.Config.GetBindAddress(), r)
}
