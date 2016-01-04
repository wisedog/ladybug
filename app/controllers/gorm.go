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
	Db, err = gorm.Open("postgres", "user=ladybug password=a1234567! dbname=ladybug sslmode=disable")
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

	// Create dummy users
	Db.AutoMigrate(&models.User{})

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)

	demoUser := &models.User{
		Name: "Demo User", Email: "demo@demo.com", Password: "demo", 
		HashedPassword: bcryptPassword, Language: "en", Region: "US",
		LastLoginAt : time.Now(), Roles : models.ROLE_ADMIN,
	}
	Db.NewRecord(demoUser) // => returns `true` if primary key is blank
	Db.Create(&demoUser)

	demoUser1 := &models.User{Name: "Wisedog", Email: "wisedog@demo.com", Password: "demo",
		HashedPassword: bcryptPassword, Language: "en", Region: "US", LastLoginAt : time.Now(),
		Roles : models.ROLE_MANAGER,
		
	}
	Db.NewRecord(demoUser1)
	Db.Create(&demoUser1)

	//Db.Model(tab).AddUniqueIndex("idx_user__gmail", "gmail")
	//Db.Model(tab).AddUniqueIndex("idx_user__pu_mail", "pu_mail")
	
	// Create dummy project
	Db.AutoMigrate(&models.Project{})

	prj := models.Project{
			Name:        "Koblentz",
			Status:      0,
			Description: "A test project",
			Prefix:      "tc",
			Users:       []models.User{*demoUser, *demoUser1},
	}
	prj1 := models.Project{
			Name:        "bremen",
			Status:      0,
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
			Prefix: "wise", Seq: 1, 
			Title: "Do not go gentle", Status: 0, Description: "Desc", SectionID: 2,
			ProjectID : prj.ID, Priority : models.PRIORITY_HIGH,
		},
		&models.TestCase{
			Prefix: "wise", Seq: 2, Title: "The Mars rover should be tested", 
			Status: 0, Description: "Desc", SectionID: 3,
			ProjectID : prj.ID, Priority : models.PRIORITY_HIGHEST,
		},
		&models.TestCase{
			Prefix: "wise", Seq: 3, Title: "I'm still arive!", 
			Status: 0, Description: "Desc", SectionID: 4,
			ProjectID : prj.ID,Priority : models.PRIORITY_MEDIUM,
		},
		&models.TestCase{
			Prefix: "wise", Seq: 4, Status: 0, SectionID: 4,
			Title: "Copy operation should be supported", 
			Description: "Copy operation is essential feature for text editing.",
			Precondition : "None",
			Steps : "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected : "Selected text are copied in TextArea" ,
			ProjectID : prj.ID, Priority : models.PRIORITY_HIGH,
		},
		&models.TestCase{
			Prefix: "wise", Seq: 5, Status: 0, SectionID: 4,
			Title: "Paste operation should be supported", 
			Description: "Paste operation is essential feature for text editing.",
			Precondition : "None",
			Steps : "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected : "Selected text are copied in TextArea",
			ProjectID : prj.ID, Priority : models.PRIORITY_LOW,
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
	
	// Create dummy build items
	Db.AutoMigrate(&models.BuildItem{})
	
	// Create dummy Review
	Db.AutoMigrate(&models.Review{})
	
	
	// Create dummy section
	Db.AutoMigrate(&models.Section{})
	sections := []*models.Section{
		&models.Section{Seq: 1, Title: "Coding Conventions", Status: 0, RootNode: true, Prefix: "wise", ProjectID : prj.ID},
		&models.Section{Seq: 1, Title: "Theme and Design", Status: 0, RootNode: true, Prefix: "wise", ProjectID : prj.ID},
		&models.Section{Seq: 1, Title: "Source code control", Status: 0, RootNode: true, Prefix: "wise", ProjectID : prj.ID},
		&models.Section{Seq: 1, Title: "Go language", RootNode: false, ParentsID: 1, Prefix: "too", ProjectID : prj.ID},
		&models.Section{Seq: 2, Title: "Javascript", RootNode: false, ParentsID: 1, Prefix: "too", ProjectID : prj.ID},
		&models.Section{Seq: 1, Title: "SB Admin2", RootNode: false, ParentsID: 2, Prefix: "aaa", ProjectID : prj.ID},
		&models.Section{Seq: 1, Title: "Git", RootNode: false, ParentsID: 3, Prefix: "tpp", ProjectID : prj.ID},
	}

	for _, section := range sections {
		Db.NewRecord(section)
		Db.Create(&section)
	}

}

func (c *GormController) Begin() revel.Result {
	txn := Db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Tx = txn
	revel.INFO.Println("c.Tx init", c.Tx)
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
	revel.INFO.Println("c.Tx commited (nil)")
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
