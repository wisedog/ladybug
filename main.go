package main

import (
  "net/http"
  "html/template"

  "github.com/gorilla/mux"
  "github.com/gorilla/schema"
  "github.com/gorilla/sessions"
  "github.com/wisedog/ladybug/controllers"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
  "github.com/wisedog/ladybug/database"
  "github.com/wisedog/ladybug/models"
  
  log "gopkg.in/inconshreveable/log15.v2"
)

// connected is private utilitiy function for checking 
// this user id now on connected or not.
func connected2(c *interfacer.AppContext, r *http.Request) *models.User{
  session, err := c.Store.Get(r, "ladybug")
  if err != nil {
    log.Warn("An error in connected(),", "msg", err)
    return nil
  }
  email := session.Values["user"]
  uid := session.Values["uid"]

  var u models.User
  c.Db.Where("email = ? and id = ?", email, uid ).First(&u)

  if u.ID == 0{
    return nil
  }
  return &u
}

type handlerFunc func(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error

//
func myHandler(c *interfacer.AppContext, h handlerFunc) http.HandlerFunc{
 return func(w http.ResponseWriter, r *http.Request) {
    errors.HandleError(w, r, func() error {
      return h(c, w, r)
    }())
  }
}


// authHandler checks authorization from session
func authHandler(c *interfacer.AppContext, h handlerFunc) http.HandlerFunc{
 return func(w http.ResponseWriter, r *http.Request) {
    errors.HandleError(w, r, func() error {
      var user *models.User
      if user = connected2(c, r); user == nil{
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
// see views/errors/404.html
func notFound (w http.ResponseWriter, r *http.Request) {
  log.Warn("Router", "404", r.URL.Path)
  t, err := template.ParseFiles("views/errors/404.html")
  if err != nil{
    log.Warn("Router", "Template ParsesFiles is failed in notFound", "msg", err)
  }
  
  t.Execute(w, nil)
}

// main is entry point of this application
func main() {
	
  log.Info("App initialize...")

  // create router with gorilla/mux
  r := mux.NewRouter()


  // load i18n resources
  controllers.LoadI18nMessage()

  ctx := &interfacer.AppContext{}
  log.Info("Initialize Database...")
  if err := database.InitDB(); err != nil{
    log.Crit("Database initialization is failed.")
    return
  }

  ctx.Db = &database.Database
  ctx.Store = sessions.NewCookieStore([]byte("ladybug"))
  ctx.Decoder = schema.NewDecoder()

  //// Routes consist of a path and a handler function.

  // application overall
  r.HandleFunc("/", myHandler(ctx, controllers.LoginPage)).Methods("GET")
  r.HandleFunc("/login", myHandler(ctx, controllers.LoginPage)).Methods("GET")
  r.HandleFunc("/login", myHandler(ctx, controllers.Login)).Methods("POST")
  r.HandleFunc("/logout", myHandler(ctx, controllers.LogOut)).Methods("GET")
  r.HandleFunc("/hello", authHandler(ctx, controllers.Welcome)).Methods("GET")

  // define user subrouter
  user := r.PathPrefix("/user/").Subrouter()
  user.HandleFunc("/profile/{id:[0-9]+}", authHandler(ctx, controllers.UserProfile)).Methods("GET")
  user.HandleFunc("/general/{id:[0-9]+}", authHandler(ctx, controllers.UserGeneral)).Methods("GET")

  // define project subrouter
  project := r.PathPrefix("/project/").Subrouter()

  // project specific
  project.HandleFunc("/{projectName}", authHandler(ctx, controllers.Dashboard)).Methods("GET")
  project.HandleFunc("/{projectName}/dashboard", authHandler(ctx, controllers.Dashboard)).Methods("GET")
  project.HandleFunc("/{projectName}/design", authHandler(ctx, controllers.DesignIndex)).Methods("GET")
  project.HandleFunc("/{projectName}/build", authHandler(ctx, controllers.BuildsIndex)).Methods("GET")
  project.HandleFunc("/{projectName}/spec", authHandler(ctx, controllers.SpecIndex)).Methods("GET")
  project.HandleFunc("/{projectName}/testplan", authHandler(ctx, controllers.PlanIndex)).Methods("GET")
  project.HandleFunc("/{projectName}/exec", authHandler(ctx, controllers.ExecIndex)).Methods("GET")

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

  // specification
  spec := project.PathPrefix("/{projectName}/spec").Subrouter()
  spec.HandleFunc("/", authHandler(ctx, controllers.SpecIndex)).Methods("GET")
  spec.HandleFunc("/list/{id:[0-9]+}", authHandler(ctx, controllers.SpecList)).Methods("GET")

  // testexec

  exec := project.PathPrefix("/{projectName}/exec").Subrouter()
  exec.HandleFunc("/", authHandler(ctx, controllers.ExecIndex)).Methods("GET")
  exec.HandleFunc("/run/{id:[0-9]+}", authHandler(ctx, controllers.ExecRun)).Methods("GET")
  exec.HandleFunc("/remove", authHandler(ctx, controllers.ExecRemove)).Methods("POST")
  exec.HandleFunc("/deny", authHandler(ctx, controllers.ExecDeny)).Methods("POST")
  exec.HandleFunc("/update", authHandler(ctx, controllers.ExecUpdateResult)).Methods("POST")
  exec.HandleFunc("/done", authHandler(ctx, controllers.ExecUpdateResult)).Methods("POST")


  r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", 
      http.FileServer(http.Dir("public/"))))
  r.NotFoundHandler = http.HandlerFunc(notFound)

  // Bind to a port and pass our router in
  http.ListenAndServe(":8000", r)
}
