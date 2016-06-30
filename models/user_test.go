package models

import(
	"testing"
)

func TestUserValidation(t *testing.T) {
  // test data is defined at database/gorm.go
  user := User{Name:"", Email:"aaa@bbb.com"}
  errorMap := user.Validate()
  
  if len(errorMap) == 0{
    t.Errorf(`User{Name:"", Email:"aaa@bbb.com"}`)
  }
  
  if errorMap["Name"] != "Name is required." {
    t.Errorf(`(User{Name:"", Email:"aaa@bbb.com"}) errorMap['name'] = %q expected %q`, errorMap["Name"], "Name is required.") 
  }
  
  user = User{Name : "Brünhild", Email:""}
  errorMap = user.Validate()
  
  if len(errorMap) == 0{
    t.Error(`User{Name : "Brünhild", Email:""}`)
  }
  
  if errorMap["Email"] != "Email is required." {
    t.Errorf(`User{Name : "Brünhild", Email:""} errorMap['Email'] = %q expected %q `, errorMap["Email"], "Email is required.")
  }
  
  user.Email = "asdadsasda"
  errorMap = user.Validate()
  if len(errorMap) == 0{
    t.Error(`User{Name : "Brünhild", Email:"asdasdad"})`)
  }
  
  if errorMap["Email"] != "Invalid email address"{
    t.Errorf(`User{Name : "Brünhild", Email:"asdadsasda"} errorMap['Email'] = %q expected %q `, errorMap["Email"], "Invalid email address")
  }
  
   
  // check empty project name
  /*
  // normal condition
  prj = Project{Name : "Lübeck", Prefix:"HL"}
  errorMap = prj.Validate()

  if len(errorMap) > 0{
    t.Errorf(`models.Project.Validate() error ({Name : "Lübeck", Prefix:"HL"})`)
  }

  if errorMap["Name"] != "" {
    t.Errorf(`({Name : "Lübeck", Prefix:"HL"}) errorMap['name'] = %q expected %q`, errorMap["Name"], "") 
  }

  // check empty prefix
  prj = Project{Name : "Hamburg", Prefix:""}
  errorMap = prj.Validate()

  if len(errorMap) == 0{
    t.Errorf(`models.Project.Validate() error ({Name : "Hamburg", Prefix:""})`)
  }

  if errorMap["Prefix"] != "Prefix is required." {
    t.Errorf(`({Name : "Hamburg", Prefix:""}) errorMap["Prefix"] = %q expected %q`, errorMap["Prefix"], "Prefix is required.") 
  }


  // check max size of prefix
    // check empty prefix
  prj = Project{Name : "Kiel", Prefix:"kielkielkielkiel"}
  errorMap = prj.Validate()

  if len(errorMap) == 0{
    t.Errorf(`models.Project.Validate() error ({Name : "Kiel", Prefix:"kielkielkielkiel"})`)
  }

  if errorMap["Prefix"] != " Size of Prefix exceeds 12 characters." {
    t.Errorf(`({Name : "Kiel", Prefix:"kielkielkielkiel"}) errorMap["Prefix"] = %q expected %q`, errorMap["Prefix"], " Size of Prefix exceeds 12 characters.") 
  }*/
}
