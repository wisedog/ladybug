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
}
