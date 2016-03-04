package main

import (
  "net/http"
  "html/template"

  "github.com/gorilla/mux"
  //"github.com/jinzhu/gorm"
  "github.com/wisedog/ladybug/controllers"
  "github.com/wisedog/ladybug/database"
  "github.com/wisedog/ladybug/interfacer"
  "github.com/wisedog/ladybug/errors"
  "github.com/gorilla/sessions"
  log "gopkg.in/inconshreveable/log15.v2"
  //"github.com/gorilla/securecookie"
)

type handlerFunc func(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error

func myHandler(c *interfacer.AppContext, h handlerFunc) http.HandlerFunc{
 return func(w http.ResponseWriter, r *http.Request) {
    errors.HandleError(w, r, func() error {
      return h(c, w, r)
    }())
  }
}

// notFound handles 404 error and render custom page. 
// see views/404.html
func notFound (w http.ResponseWriter, r *http.Request) {
  t, err := template.ParseFiles("views/404.html")
  if err != nil{
    log.Warn("Template ParsesFiles is failed in notFound", "msg", err)
  }
  
  t.Execute(w, nil)
}

// main is entry point of this application
func main() {
	
  log.Info("App initialize...")

  // create router with gorilla/mux
  r := mux.NewRouter()

  ctx := &interfacer.AppContext{}
  log.Info("Initialize Database...")
  if err := database.InitDB(); err != nil{
    log.Crit("Database initialization is failed.")
    return
  }

  ctx.Db = &database.Database
  ctx.Store = sessions.NewCookieStore([]byte("ladybug"))

  

  //// Routes consist of a path and a handler function.

  // application overall
  r.HandleFunc("/", myHandler(ctx, controllers.LoginPage)).Methods("GET")
  r.HandleFunc("/login", myHandler(ctx, controllers.LoginPage)).Methods("GET")
  r.HandleFunc("/login", myHandler(ctx, controllers.Login)).Methods("POST")
  r.HandleFunc("/logout", myHandler(ctx, controllers.LogOut)).Methods("GET")
  r.HandleFunc("/hello", myHandler(ctx, controllers.Welcome)).Methods("GET")

  // define user subrouter
  user := r.PathPrefix("/user/").Subrouter()
  user.HandleFunc("/profile/{id:[0-9]+}", myHandler(ctx, controllers.UserProfile)).Methods("GET")
  user.HandleFunc("/general/{id:[0-9]+}", myHandler(ctx, controllers.UserGeneral)).Methods("GET")

  // define project subrouter
  project := r.PathPrefix("/project/").Subrouter()
  project.HandleFunc("/{projectName}", myHandler(ctx, controllers.Dashboard)).Methods("GET")
  project.HandleFunc("/{projectName}/dashboard", myHandler(ctx, controllers.Dashboard)).Methods("GET")

  // define test design
  project.HandleFunc("/{projectName}/design", myHandler(ctx, controllers.DesignIndex)).Methods("GET")

  // section
  section := project.PathPrefix("/{projectName}/section/").Subrouter()
  section.HandleFunc("/testcase/{sectionID}", myHandler(ctx, controllers.GetAllTestCases)).Methods("GET")
  section.HandleFunc("/add/{sectionID}", myHandler(ctx, controllers.SectionAdd)).Methods("GET")

  /*POST    /project/:project/case/save             TestCases.Save
  GET     /project/:project/case/:id/edit         TestCases.Edit
  POST    /project/:project/case/:id/delete       TestCases.Delete
  GET     /project/:project/case/:id              TestCases.CaseIndex*/
  // test cases
  testcase := project.PathPrefix("/{projectName}/case/").Subrouter()
  testcase.HandleFunc("/view/{id:[0-9]+}", myHandler(ctx, controllers.CaseView)).Methods("GET")
  //testcase.HandleFunc("/save", myHandler(ctx, controllers.CaseSave)).Methods("POST")
  //testcase.HandleFunc("/edit/{id:[0-9]+", myHandler(ctx, controllers.CaseEdit)).Methods("GET")
  //testcase.HandleFunc("/add", myHandler(ctx, controllers.CaseAdd)).Methods("GET")


  // test plan


  // builds

  // testcase

  


  //r.PathPrefix("/public/").Handler(http.FileServer(http.Dir("./public/")))
  r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", 
      http.FileServer(http.Dir("public/"))))
  r.NotFoundHandler = http.HandlerFunc(notFound)

  // Bind to a port and pass our router in
  http.ListenAndServe(":8000", r)
}