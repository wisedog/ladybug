package controllers

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type GormController struct {
	*revel.Controller
	Tx *gorm.DB
}

var Db gorm.DB

func InitDB() {
	var err error
	Db, err = gorm.Open("postgres", "user=postgres dbname=ladybug sslmode=disable")
	if err != nil {
		revel.ERROR.Println("FATAL", err)
		panic(err)
	}

	if revel.Config.BoolDefault("Ladybug.droptable", true) == false {
		Db.AutoMigrate(&models.User{})
		Db.AutoMigrate(&models.TestCase{})
		Db.AutoMigrate(&models.TestPlan{})
		Db.AutoMigrate(&models.Build{})
		Db.AutoMigrate(&models.BuildItem{})
		Db.AutoMigrate(&models.Section{})
		Db.AutoMigrate(&models.Execution{})
		Db.AutoMigrate(&models.TestResult{})
		Db.AutoMigrate(&models.Review{})
		Db.AutoMigrate(&models.Category{})
		Db.AutoMigrate(&models.Specification{})

		revel.INFO.Println("All tables are not dropped!")
	} else {
		revel.INFO.Println("All tables are DROPPED!")
		createDummy()
	}

}

func createDummy() {

	// drop all table while on development phase
	Db.DropTable(&models.User{})
	Db.DropTable(&models.Project{})
	Db.DropTable(&models.TestCase{})
	Db.DropTable(&models.TestPlan{})
	Db.DropTable(&models.Build{})
	Db.DropTable(&models.BuildItem{})
	Db.DropTable(&models.Section{})
	Db.DropTable(&models.Execution{})
	Db.DropTable(&models.Review{})
	Db.DropTable(&models.TestResult{})
	Db.DropTable(&models.Category{})
	Db.DropTable(&models.Specification{})
	Db.DropTable(&models.Activity{})
	Db.DropTable(&models.History{})

	// Create dummy users
	Db.AutoMigrate(&models.User{})

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)

	demoUser := &models.User{
		Name: "Rey", Email: "demo@demo.com", Password: "demo", 
		HashedPassword: bcryptPassword, Language: "en", Region: "US",
		LastLoginAt : time.Now(), Roles : models.ROLE_ADMIN,
		Photo : "rey_160x160", Location : "Jakku",
		Notes : "I know all about waiting. For my family. They'll be back, one day.",
	}
	Db.NewRecord(demoUser) // => returns `true` if primary key is blank
	Db.Create(&demoUser)

	demoUser1 := &models.User{Name: "Poe Dameron", Email: "wisedog@demo.com", Password: "demo",
		HashedPassword: bcryptPassword, Language: "en", Region: "US", LastLoginAt : time.Now(),
		Roles : models.ROLE_MANAGER, Photo : "poe_160x160",
		Location : "D'Qar", Notes : "Red squad, blue squad, take my lead.",
		
	}
	Db.NewRecord(demoUser1)
	Db.Create(&demoUser1)

	//Db.Model(tab).AddUniqueIndex("idx_user__gmail", "gmail")
	//Db.Model(tab).AddUniqueIndex("idx_user__pu_mail", "pu_mail")
	
	// Create dummy project
	Db.AutoMigrate(&models.Project{})

	prj := models.Project{
			Name:        "Koblentz",
			Status:      1,
			Description: "A test project",
			Prefix:      "tc",
			Users:       []models.User{*demoUser, *demoUser1},
	}
	prj1 := models.Project{
			Name:        "bremen",
			Status:      1,
			Description: "A test project2",
			Prefix:      "wise",
			Users:       []models.User{*demoUser, *demoUser1},
	}

	Db.NewRecord(prj)
	Db.Create(&prj)
	Db.NewRecord(prj1)
	Db.Create(&prj1)

	// Create dummy testcases
	Db.AutoMigrate(&models.TestCase{})

/*
1. Drag some texts on web browser to select text and CTRL + C 
2. Click TextArea and CTRL + V
*/
	testcases := []*models.TestCase{
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 1, 
			Title: "Do not go gentle", Status: models.TC_STATUS_ACTIVATE, Description: "Desc", SectionID: 2,
			ProjectID : prj.ID, Priority : models.PRIORITY_HIGH, CategoryID : 1,
			DisplayID : prj.Prefix + "-1", Estimated : 10,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 2, Title: "The Mars rover should be tested", 
			Status: models.TC_STATUS_ACTIVATE, Description: "Desc", SectionID: 3,
			ProjectID : prj.ID, Priority : models.PRIORITY_HIGHEST,
			CategoryID : 1, DisplayID : prj.Prefix + "-2",Estimated : 1,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 3, Title: "I'm still arive!",  CategoryID : 3,
			Status: models.TC_STATUS_ACTIVATE, Description: "Desc", SectionID: 4,
			ProjectID : prj.ID,Priority : models.PRIORITY_MEDIUM,
			DisplayID : prj.Prefix + "-3",Estimated : 4,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 4, Status: models.TC_STATUS_DRAFT, SectionID: 4, CategoryID : 2,
			Title: "Copy operation should be supported", 
			Description: "Copy operation is essential feature for text editing.",
			Precondition : "None",
			Steps : "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected : "Selected text are copied in TextArea" ,
			ProjectID : prj.ID, Priority : models.PRIORITY_HIGH,
			DisplayID : prj.Prefix + "-4",Estimated : 5,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 5, Status: models.TC_STATUS_INACTIVE, SectionID: 4, CategoryID : 1,
			Title: "Paste operation should be supported", 
			Description: "Paste operation is essential feature for text editing.",
			Precondition : "None",
			Steps : "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected : "Selected text are copied in TextArea",
			ProjectID : prj.ID, Priority : models.PRIORITY_LOW,
			DisplayID : prj.Prefix + "-5",
			},
	}

	for _, tc := range testcases {
		Db.NewRecord(tc)
		Db.Create(&tc)
	}

	// Create dummy testplan
	Db.AutoMigrate(&models.TestPlan{})
	
	// Create dummy test execution
	Db.AutoMigrate(&models.Execution{})
	
	// Create dummy test result
	Db.AutoMigrate(&models.TestResult{})

	// Create dummy build
	Db.AutoMigrate(&models.Build{})
	b := []*models.Build{
		&models.Build{Name:"Millenium Falcon", 
			Description : "Modeling files for Millenium Falcon", 
			Project_id : prj.ID, ToolName: "manual",
		},
	}
	
	for _, bi := range b {
		Db.NewRecord(bi)
		Db.Create(&bi)
	}
	
	
	// Create dummy build items
	Db.AutoMigrate(&models.BuildItem{})
	
	// Create dummy Review
	Db.AutoMigrate(&models.Review{})
	
	
	// Create dummy section
	Db.AutoMigrate(&models.Section{})
	sections := []*models.Section{
		&models.Section{Seq: 1, Title: "Coding Conventions", Status: 0, RootNode: true, 
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "Theme and Design", Status: 0, RootNode: true, 
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "Source code control", Status: 0, RootNode: true, 
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "Go language", RootNode: false, 
			ParentsID: 1, Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 2, Title: "Javascript", RootNode: false, 
			ParentsID: 1, Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "SB Admin2", RootNode: false, ParentsID: 2, 
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "Git", RootNode: false, ParentsID: 3, 
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "AAA", RootNode: true,  
			Prefix: prj1.Prefix, ProjectID : prj1.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "BBB", RootNode: false, ParentsID : 8,  
			Prefix: prj1.Prefix, ProjectID : prj1.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "ccc", RootNode: true,  
			Prefix: prj1.Prefix, ProjectID : prj1.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "ddd", RootNode: false, ParentsID : 10, 
			Prefix: prj1.Prefix, ProjectID : prj1.ID, ForTestCase : true,},
		&models.Section{Seq: 1, Title: "Test Specifications", RootNode: true,
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : false,},
		&models.Section{Seq: 1, Title: "Functional", RootNode: false,
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : false,
			ParentsID : 12,
		},
		&models.Section{Seq: 1, Title: "Non-Functional", RootNode: false,
			Prefix: prj.Prefix, ProjectID : prj.ID, ForTestCase : false,
			ParentsID : 12,
		},
	}

	for _, section := range sections {
		Db.NewRecord(section)
		Db.Create(&section)
	}
	
	Db.AutoMigrate(&models.Category{})
	cate := []*models.Category{
		&models.Category{Name : "Funtionality"},
		&models.Category{Name : "Performance"},
		&models.Category{Name : "Usability"},
		&models.Category{Name : "Regression"},
		&models.Category{Name : "Automated"},
		&models.Category{Name : "Security"},
		&models.Category{Name : "Compatibility"},
		&models.Category{Name : "Accessability"},
	}
	
	for _, ct := range cate {
		Db.NewRecord(ct)
		Db.Create(&ct)
	}
	
	Db.AutoMigrate(&models.Specification{})
	
	specs := []*models.Specification{
		&models.Specification{Name:"Hello, world", SectionID : 13, 
			Status : models.SPEC_STATUS_ACTIVATE, Priority: models.PRIORITY_HIGH,
		},
		&models.Specification{Name:"Hello, stranger", SectionID : 14,
			Status : models.SPEC_STATUS_ACTIVATE, Priority: models.PRIORITY_MEDIUM,
		},
		&models.Specification{Name:"Good bye", SectionID : 14,
			Status : models.SPEC_STATUS_ACTIVATE, Priority: models.PRIORITY_LOW,
		},
	}
	
	for _, sp := range specs {
		Db.NewRecord(sp)
		Db.Create(&sp)
	}
	
	
	// for creating dummy for Activity
	Db.AutoMigrate(&models.Activity{})
	
	activities := []*models.Activity{
		&models.Activity{UserID : demoUser.ID, Content : "Rey finished Test Execution #1"},
		&models.Activity{UserID : demoUser.ID, Content : "Rey created Test Plan #1"},
		&models.Activity{UserID : demoUser.ID, Content : "Rey modified TC-5"},
	}
	
	for _, ac := range activities {
		Db.NewRecord(ac)
		Db.Create(&ac)
	}
	
	// for creating dummy for History
	Db.AutoMigrate(&models.History{})
}

func (c *GormController) Begin() revel.Result {
	txn := Db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Tx = txn
	//revel.INFO.Println("c.Tx init", c.Tx)
	return nil
}

func (c *GormController) Commit() revel.Result {
	if c.Tx == nil {
		return nil
	}
	c.Tx.Commit()
	if err := c.Tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Tx = nil
	//revel.INFO.Println("c.Tx commited (nil)")
	return nil
}

func (c *GormController) Rollback() revel.Result {
	if c.Tx == nil {
		return nil
	}
	c.Tx.Rollback()
	if err := c.Tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Tx = nil
	return nil
}
