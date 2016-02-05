package controllers

import (
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(TestPlans.checkUser, revel.BEFORE)
	//revel.InterceptMethod(TestCases.checkUser, revel.BEFORE)
	//	revel.InterceptFunc(checkUser, revel.BEFORE, &Application{})
	revel.InterceptMethod((*GormController).Commit, revel.AFTER)
	revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)
	
	revel.TemplateFuncs["priority2str"] = func(p int) string { 
		str := ""
		switch p{
			case 1: str = "Highest"	// FIXME I want to use controller.Message... 
			case 2: str = "High"
			case 3: str = "Medium"
			case 4: str = "Low"
			case 5: str = "Lowest"
			default: str = "unknown"
		}
		return str
	}
}