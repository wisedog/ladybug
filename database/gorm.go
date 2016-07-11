package database

import (
  "fmt"
  "time"

  "golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
	
  log "gopkg.in/inconshreveable/log15.v2" 
)

var Database *gorm.DB

// InitDB initialize the database and create dummies if it needs
func InitDB(conf *interfacer.AppConfig) (*gorm.DB, error){

  args := getDialectArgs(conf)
  var err error
	//Database, err = gorm.Open("postgres", "user=ladybug dbname=ladybug port=5432 sslmode=disable")
  Database, err = gorm.Open("postgres", args)
	if err != nil {
		log.Info("Database", "msg", err.Error())
		return Database, err
	}

  Database.AutoMigrate(&models.User{})
  Database.AutoMigrate(&models.TestCase{})
  Database.AutoMigrate(&models.TestPlan{})
  Database.AutoMigrate(&models.Build{})
  Database.AutoMigrate(&models.BuildItem{})
  Database.AutoMigrate(&models.Section{})
  Database.AutoMigrate(&models.Execution{})
  Database.AutoMigrate(&models.TestResult{})
  Database.AutoMigrate(&models.Review{})
  Database.AutoMigrate(&models.Category{})
  Database.AutoMigrate(&models.Specification{})
	Database.AutoMigrate(&models.Milestone{})
  createDummy()

  return Database, nil
}

// getDialectArgs returns argument of dialect database. 
func getDialectArgs(conf *interfacer.AppConfig) string{
  driver := conf.GetValue("db.driver")

  username := conf.GetValue("db.username")
  pwd := conf.GetValue("db.password")
  port := conf.GetValue("db.port")
  address := conf.GetValue("db.address")
  dbname := conf.GetValue("db.database_name")
  //extraParam := conf.GetValue("db.extra_param")

  var args string
  switch driver{
    case "postgres":
      args = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, pwd, address, port, dbname)
    case "mysql":
      args="" //TODO

    default:
      args = ""
  }

  return args
}


// createDummy creates dummy data for database
func createDummy() {

	// drop all table while on development phase
	Database.DropTable(&models.User{})
	Database.DropTable(&models.Project{})
	Database.DropTable(&models.TestCase{})
	Database.DropTable(&models.TestPlan{})
	Database.DropTable(&models.Build{})
	Database.DropTable(&models.BuildItem{})
	Database.DropTable(&models.Section{})
	Database.DropTable(&models.Execution{})
	Database.DropTable(&models.Review{})
	Database.DropTable(&models.TestResult{})
	Database.DropTable(&models.Category{})
	Database.DropTable(&models.Specification{})
	Database.DropTable(&models.Activity{})
	Database.DropTable(&models.History{})
	Database.DropTable(&models.Milestone{})

	// Create dummy users
	Database.AutoMigrate(&models.User{})

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)

	demoUser := &models.User{
		Name: "Rey", Email: "demo@demo.com", Password: "demo", 
		HashedPassword: bcryptPassword, Language: "en", Region: "US",
		LastLoginAt : time.Now(), Roles : models.RoleAdmin,
		Photo : "rey_160x160", Location : "Jakku",
		Notes : "I know all about waiting. For my family. They'll be back, one day.",
	}
	Database.NewRecord(demoUser) // => returns `true` if primary key is blank
	Database.Create(&demoUser)

	demoUser1 := &models.User{Name: "Poe Dameron", Email: "wisedog@demo.com", Password: "demo",
		HashedPassword: bcryptPassword, Language: "en", Region: "US", LastLoginAt : time.Now(),
		Roles : models.RoleManager, Photo : "poe_160x160",
		Location : "D'Qar", Notes : "Red squad, blue squad, take my lead.",
		
	}
	Database.NewRecord(demoUser1)
	Database.Create(&demoUser1)

	//Database.Model(tab).AddUniqueIndex("idx_user__gmail", "gmail")
	//Database.Model(tab).AddUniqueIndex("idx_user__pu_mail", "pu_mail")
	
	// Create dummy project
	Database.AutoMigrate(&models.Project{})

	prj := models.Project{
			Name:        "Koblentz",
			Status:      1,
			Description: "A test project",
			Prefix:      "TC",
			Users:       []models.User{*demoUser, *demoUser1},
	}
	prj1 := models.Project{
			Name:        "bremen",
			Status:      1,
			Description: "A test project2",
			Prefix:      "wise",
			Users:       []models.User{*demoUser, *demoUser1},
	}

	Database.NewRecord(prj)
	Database.Create(&prj)
	Database.NewRecord(prj1)
	Database.Create(&prj1)

	// Create dummy testcases
	Database.AutoMigrate(&models.TestCase{})
  
	testcases := []*models.TestCase{
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 1, 
			Title: "Do not go gentle", Status: models.TcStatusActivate, Description: "Desc", SectionID: 2,
			ProjectID : prj.ID, Priority : models.PriorityHigh, CategoryID : 1,
			DisplayID : prj.Prefix + "-1", Estimated : 10,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 2, Title: "The Mars rover should be tested", 
			Status: models.TcStatusActivate, Description: "Desc", SectionID: 3,
			ProjectID : prj.
      ID, Priority : models.PriorityHighest,
			CategoryID : 1, DisplayID : prj.Prefix + "-2",Estimated : 1,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 3, Title: "I'm still arive!",  CategoryID : 3,
			Status: models.TcStatusActivate, Description: "Desc", SectionID: 4,
			ProjectID : prj.ID,Priority : models.PriorityMedium,
			DisplayID : prj.Prefix + "-3",Estimated : 4,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 4, Status: models.TcStatusActivate, SectionID: 4, CategoryID : 2,
			Title: "Copy operation should be supported", 
			Description: "Copy operation is essential feature for text editing.",
			Precondition : "None",
			Steps : "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected : "Selected text are copied in TextArea" ,
			ProjectID : prj.ID, Priority : models.PriorityHigh,
			DisplayID : prj.Prefix + "-4",Estimated : 5,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 5, Status: models.TcStatusActivate, SectionID: 4, CategoryID : 1,
			Title: "Paste operation should be supported", 
			Description: "Paste operation is essential feature for text editing.",
			Precondition : "None",
			Steps : "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected : "Selected text are copied in TextArea",
			ProjectID : prj.ID, Priority : models.PriorityLow,
			DisplayID : prj.Prefix + "-5",
			},
	}

	for _, tc := range testcases {
		Database.NewRecord(tc)
		Database.Create(&tc)
	}
  
  /*var tempTestCase []models.TestCase
  Database.Find(&tempTestCase)
  
  now := time.Now()
  for i, iter := range tempTestCase{
    then := now.AddDate(0, 0, -i)
    fmt.Println("then:", then)
    tc.CreatedAt = then
  }*/

	// Create dummy testplan
	Database.AutoMigrate(&models.TestPlan{})
	
	// Create dummy test execution
	Database.AutoMigrate(&models.Execution{})
	
	// Create dummy test result
	Database.AutoMigrate(&models.TestResult{})

	// Create dummy build
	Database.AutoMigrate(&models.Build{})
	b := []*models.Build{
		&models.Build{Name:"Millenium Falcon", 
			Description : "Modeling files for Millenium Falcon", 
			Project_id : prj.ID, ToolName: "manual",
		},
	}
	
	for _, bi := range b {
		Database.NewRecord(bi)
		Database.Create(&bi)
	}
	
	
	// Create dummy build items
	Database.AutoMigrate(&models.BuildItem{})
	
	// Create dummy Review
	Database.AutoMigrate(&models.Review{})
	
	
	// Create dummy section
	Database.AutoMigrate(&models.Section{})
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
		Database.NewRecord(section)
		Database.Create(&section)
	}
	
	Database.AutoMigrate(&models.Category{})
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
		Database.NewRecord(ct)
		Database.Create(&ct)
	}
	
	Database.AutoMigrate(&models.Specification{})
	
	specs := []*models.Specification{
		&models.Specification{Name:"Hello, world", SectionID : 13, 
			Status : models.SpecStatusActivate, Priority: models.PriorityHigh,
		},
		&models.Specification{Name:"Hello, stranger", SectionID : 14,
			Status : models.SpecStatusActivate, Priority: models.PriorityMedium,
		},
		&models.Specification{Name:"Good bye", SectionID : 14,
			Status : models.SpecStatusActivate, Priority: models.PriorityLow,
		},
	}
	
	for _, sp := range specs {
		Database.NewRecord(sp)
		Database.Create(&sp)
	}
	
	
	// for creating dummy for Activity
	Database.AutoMigrate(&models.Activity{})
	
	activities := []*models.Activity{
		&models.Activity{UserID : demoUser.ID, Content : "Rey finished Test Execution #1"},
		&models.Activity{UserID : demoUser.ID, Content : "Rey created Test Plan #1"},
		&models.Activity{UserID : demoUser.ID, Content : "Rey modified TC-5"},
	}
	
	for _, ac := range activities {
		Database.NewRecord(ac)
		Database.Create(&ac)
	}
	
	// for creating dummy for History
	Database.AutoMigrate(&models.History{})

	Database.AutoMigrate(&models.Milestone{})
	now := time.Now()
	oneMonthFromNow := time.Hour * 24 * 30

	next := now.Add(oneMonthFromNow)
	next1 := next.Add(oneMonthFromNow)
	Milestones := []*models.Milestone{
		&models.Milestone{Name : "Milestone#1", ProjectID : prj.ID, Status : models.MilestoneStatusActive, DueDate : next},
		&models.Milestone{Name : "Milestone#2", ProjectID : prj.ID, Status : models.MilestoneStatusActive, DueDate : next1},
	}

	for _, ms := range Milestones {
		Database.NewRecord(ms)
		Database.Create(&ms)
	}

}
/*
func (c *GormController) Begin() revel.Result {
	txn := Database.Begin()
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
*/