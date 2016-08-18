package database

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	// pq library is used by gorm
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/models"

	log "gopkg.in/inconshreveable/log15.v2"
)

// Database is an instance of connection to DB
var Database *gorm.DB

// InitDB initialize the database and create dummies if it needs
func InitDB(conf *interfacer.AppConfig) (*gorm.DB, error) {

	var err error
	Database, err = gorm.Open("postgres", getDialectArgs(conf))
	if err != nil {
		log.Info("Database", "msg", err.Error())
		return Database, err
	}
	//defer Database.Close()

	Database.AutoMigrate(&models.User{})
	Database.AutoMigrate(&models.TestCase{})
	Database.AutoMigrate(&models.TestPlan{})
	Database.AutoMigrate(&models.Build{})
	Database.AutoMigrate(&models.BuildItem{})
	Database.AutoMigrate(&models.Section{})
	Database.AutoMigrate(&models.Execution{})
	Database.AutoMigrate(&models.TestCaseResult{})
	Database.AutoMigrate(&models.Review{})
	Database.AutoMigrate(&models.Category{})
	Database.AutoMigrate(&models.Requirement{})
	Database.AutoMigrate(&models.Milestone{})
	Database.AutoMigrate(&models.ReqType{})
	createDummy()

	return Database, nil
}

func loadDefault() error {
	return nil
}

// getDialectArgs returns argument of dialect database.
func getDialectArgs(conf *interfacer.AppConfig) string {
	var args string
	// not set by argument
	if conf.GetDialect() == "" {
		driver := conf.GetValue("db.driver")
		username := conf.GetValue("db.username")
		pwd := conf.GetValue("db.password")
		port := conf.GetValue("db.port")
		address := conf.GetValue("db.address")
		dbname := conf.GetValue("db.database_name")
		//extraParam := conf.GetValue("db.extra_param")
		switch driver {
		case "postgres":
			args = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, pwd, address, port, dbname)
		case "mysql":
			args = "" //TODO

		default:
			args = ""
		}
	} else {
		// set by argument
		args = conf.GetDialect()
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
	Database.DropTable(&models.TestCaseResult{})
	Database.DropTable(&models.Category{})
	Database.DropTable(&models.Requirement{})
	Database.DropTable(&models.Activity{})
	Database.DropTable(&models.History{})
	Database.DropTable(&models.Milestone{})
	Database.DropTable(&models.TcReqRelationHistory{})
	Database.DropTable(&models.ReqType{})

	// Create dummy users
	Database.AutoMigrate(&models.User{})

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)

	demoUser := &models.User{
		Name: "Rey", Email: "demo@demo.com", Password: "demo",
		HashedPassword: bcryptPassword, Language: "en", Region: "US",
		LastLoginAt: time.Now(), Roles: models.RoleAdmin,
		Photo: "rey_160x160", Location: "Jakku",
		Notes: "I know all about waiting. For my family. They'll be back, one day.",
	}
	Database.NewRecord(demoUser) // => returns `true` if primary key is blank
	Database.Create(&demoUser)

	demoUser1 := &models.User{Name: "Poe Dameron", Email: "wisedog@demo.com", Password: "demo",
		HashedPassword: bcryptPassword, Language: "en", Region: "US", LastLoginAt: time.Now(),
		Roles: models.RoleManager, Photo: "poe_160x160",
		Location: "D'Qar", Notes: "Red squad, blue squad, take my lead.",
	}
	Database.NewRecord(demoUser1)
	Database.Create(&demoUser1)

	//Database.Model(tab).AddUniqueIndex("idx_user__gmail", "gmail")
	//Database.Model(tab).AddUniqueIndex("idx_user__pu_mail", "pu_mail")

	// Create dummy project
	Database.AutoMigrate(&models.Project{})

	prj := models.Project{
		Name:        "Sample Project",
		Status:      1,
		Description: "A sample project. If you are used to Ladybug, remove this project",
		Prefix:      "TC",
		Users:       []models.User{*demoUser, *demoUser1},
	}
	prj1 := models.Project{
		Name:        "Another Sample Project",
		Status:      1,
		Description: "Second sample project. If you are used to Ladybug, remove this project",
		Prefix:      "wise",
		Users:       []models.User{*demoUser, *demoUser1},
	}

	Database.NewRecord(prj)
	Database.Create(&prj)
	Database.NewRecord(prj1)
	Database.Create(&prj1)

	// Create dummy testcases
	Database.AutoMigrate(&models.TestCase{})

	now := time.Now()

	testcases := []*models.TestCase{
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 1,
			Title: "Do not go gentle", Status: models.TcStatusActivate, Description: "Desc", SectionID: 2,
			ProjectID: prj.ID, Priority: models.PriorityHigh, CategoryID: 1,
			DisplayID: prj.Prefix + "-1", Estimated: 10,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 2, Title: "The Mars rover should be tested",
			Status: models.TcStatusActivate, Description: "Desc", SectionID: 3,
			ProjectID: prj.
				ID, Priority: models.PriorityHighest,
			CategoryID: 1, DisplayID: prj.Prefix + "-2", Estimated: 1,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 3, Title: "I'm still arive!", CategoryID: 3,
			Status: models.TcStatusActivate, Description: "Desc", SectionID: 4,
			ProjectID: prj.ID, Priority: models.PriorityMedium,
			DisplayID: prj.Prefix + "-3", Estimated: 4,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 4, Status: models.TcStatusActivate, SectionID: 4, CategoryID: 2,
			Title:        "Copy operation should be supported",
			Description:  "Copy operation is essential feature for text editing.",
			Precondition: "None",
			Steps:        "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected:     "Selected text are copied in TextArea",
			ProjectID:    prj.ID, Priority: models.PriorityHigh,
			DisplayID: prj.Prefix + "-4", Estimated: 5,
		},
		&models.TestCase{
			Prefix: prj.Prefix, Seq: 5, Status: models.TcStatusActivate, SectionID: 4, CategoryID: 1,
			Title:        "Paste operation should be supported",
			Description:  "Paste operation is essential feature for text editing.",
			Precondition: "None",
			Steps:        "1. Drag some texts on web browser to select text and CTRL + C\n2. Click TextArea and CTRL + V",
			Expected:     "Selected text are copied in TextArea",
			ProjectID:    prj.ID, Priority: models.PriorityLow,
			DisplayID: prj.Prefix + "-5",
		},
	}

	for _, tc := range testcases {
		Database.NewRecord(tc)
		Database.Create(&tc)
	}

	// belows are for sample purpose
	var tempTestCase []models.TestCase
	Database.Find(&tempTestCase)

	n1 := now.AddDate(0, 0, -3)
	n2 := now.AddDate(0, 0, -8)
	n3 := now.AddDate(0, 0, -15)
	n4 := now.AddDate(0, 0, -22)
	timeArray := [4]time.Time{n1, n2, n3, n4}

	for i := 0; i < len(tempTestCase); i++ {
		k := tempTestCase[i]
		n := timeArray[rand.Intn(4)]
		Database.Model(&k).Update("created_at", n)
	}

	// Create dummy testplan
	Database.AutoMigrate(&models.TestPlan{})

	// Create dummy test execution
	Database.AutoMigrate(&models.Execution{})

	// Create dummy test result
	Database.AutoMigrate(&models.TestCaseResult{})

	// Create dummy build
	Database.AutoMigrate(&models.Build{})
	b := []*models.Build{
		&models.Build{Name: "Millenium Falcon",
			Description: "Modeling files for Millenium Falcon",
			Project_id:  prj.ID, ToolName: "manual",
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
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "Theme and Design", Status: 0, RootNode: true,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "Source code control", Status: 0, RootNode: true,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "Go language", RootNode: false,
			ParentsID: 1, Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 2, Title: "Javascript", RootNode: false,
			ParentsID: 1, Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "SB Admin2", RootNode: false, ParentsID: 2,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "Git", RootNode: false, ParentsID: 3,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "AAA", RootNode: true,
			Prefix: prj1.Prefix, ProjectID: prj1.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "BBB", RootNode: false, ParentsID: 8,
			Prefix: prj1.Prefix, ProjectID: prj1.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "ccc", RootNode: true,
			Prefix: prj1.Prefix, ProjectID: prj1.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "ddd", RootNode: false, ParentsID: 10,
			Prefix: prj1.Prefix, ProjectID: prj1.ID, ForTestCase: true},
		&models.Section{Seq: 1, Title: "Requirements", RootNode: true,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: false},
		&models.Section{Seq: 1, Title: "Functional", RootNode: false,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: false,
			ParentsID: 12,
		},
		&models.Section{Seq: 1, Title: "Non-Functional", RootNode: false,
			Prefix: prj.Prefix, ProjectID: prj.ID, ForTestCase: false,
			ParentsID: 12,
		},
	}

	for _, section := range sections {
		Database.NewRecord(section)
		Database.Create(&section)
	}

	Database.AutoMigrate(&models.Category{})
	cate := []*models.Category{
		&models.Category{Name: "Funtionality"},
		&models.Category{Name: "Performance"},
		&models.Category{Name: "Usability"},
		&models.Category{Name: "Regression"},
		&models.Category{Name: "Automated"},
		&models.Category{Name: "Security"},
		&models.Category{Name: "Compatibility"},
		&models.Category{Name: "Accessability"},
	}

	for _, ct := range cate {
		Database.NewRecord(ct)
		Database.Create(&ct)
	}

	Database.AutoMigrate(&models.Requirement{})

	reqs := []*models.Requirement{
		&models.Requirement{Title: "Hello, world", SectionID: 13,
			Description: "hello", ReqTypeID: 1,
			Status: models.ReqStatusDraft, Priority: models.PriorityHigh, ProjectID: prj.ID,
		},
		&models.Requirement{Title: "Hello, stranger", SectionID: 14,
			Description: "blahblah", ReqTypeID: 2,
			Status: models.ReqStatusDraft, Priority: models.PriorityMedium, ProjectID: prj.ID,
			RelatedTestCases: []models.TestCase{*testcases[1], *testcases[4]},
		},
		&models.Requirement{Title: "Good bye", SectionID: 14,
			Description: "aaaa", ReqTypeID: 3,
			Status: models.ReqStatusDraft, Priority: models.PriorityLow, ProjectID: prj.ID,
			RelatedTestCases: []models.TestCase{*testcases[1], *testcases[2]},
		},
	}

	for _, sp := range reqs {
		Database.NewRecord(sp)
		Database.Create(&sp)
	}

	var tmpReqs []models.Requirement
	Database.Find(&tmpReqs)

	timeArrayReq := [3]time.Time{n4, n3, n2}

	for i := 0; i < len(tmpReqs); i++ {
		k := tmpReqs[i]
		n := timeArrayReq[i]
		Database.Model(&k).Update("created_at", n)
	}

	// for creating dummy for Activity
	Database.AutoMigrate(&models.Activity{})

	activities := []*models.Activity{
		&models.Activity{UserID: demoUser.ID, Content: "Rey finished Test Execution #1"},
		&models.Activity{UserID: demoUser.ID, Content: "Rey created Test Plan #1"},
		&models.Activity{UserID: demoUser.ID, Content: "Rey modified TC-5"},
	}

	for _, ac := range activities {
		Database.NewRecord(ac)
		Database.Create(&ac)
	}

	// for creating dummy for History
	Database.AutoMigrate(&models.History{})

	Database.AutoMigrate(&models.Milestone{})
	oneMonthFromNow := time.Hour * 24 * 30

	next := now.Add(oneMonthFromNow)
	next1 := next.Add(oneMonthFromNow)
	Milestones := []*models.Milestone{
		&models.Milestone{Name: "Milestone#1", ProjectID: prj.ID, Status: models.MilestoneStatusActive, DueDate: next},
		&models.Milestone{Name: "Milestone#2", ProjectID: prj.ID, Status: models.MilestoneStatusActive, DueDate: next1},
	}

	for _, ms := range Milestones {
		Database.NewRecord(ms)
		Database.Create(&ms)
	}

	// for testcase-requirement relationship
	Database.AutoMigrate(&models.TcReqRelationHistory{})

	relations := []*models.TcReqRelationHistory{
		&models.TcReqRelationHistory{
			RequirementID: 1, TestCaseID: 2,
			Kind:      models.TcReqRelationHistoryLink,
			ProjectID: prj.ID,
		},
		&models.TcReqRelationHistory{
			RequirementID: 2, TestCaseID: 2,
			Kind:      models.TcReqRelationHistoryLink,
			ProjectID: prj.ID,
		},
		&models.TcReqRelationHistory{
			RequirementID: 2, TestCaseID: 3,
			Kind:      models.TcReqRelationHistoryUnlink,
			ProjectID: prj.ID,
		},
		&models.TcReqRelationHistory{
			RequirementID: 3, TestCaseID: 4,
			Kind:      models.TcReqRelationHistoryLink,
			ProjectID: prj.ID,
		},
	}

	for _, ac := range relations {
		Database.NewRecord(ac)
		Database.Create(&ac)
	}

	var tempRelation []models.TcReqRelationHistory
	Database.Find(&tempRelation)

	tmpTimeArray := []time.Time{n4, n3, n3, n2}
	for i := 0; i < len(tempRelation); i++ {
		k := tempRelation[i]
		n := tmpTimeArray[i]
		Database.Model(&k).Update("created_at", n)
	}

	Database.AutoMigrate(&models.ReqType{})

	reqtypes := []*models.ReqType{
		&models.ReqType{Name: "Use Case"},
		&models.ReqType{Name: "Information"},
		&models.ReqType{Name: "Feature"},
		&models.ReqType{Name: "User Interface"},
		&models.ReqType{Name: "Non Functional"},
		&models.ReqType{Name: "Constraint"},
		&models.ReqType{Name: "System Function"},
	}

	for _, t := range reqtypes {
		Database.NewRecord(t)
		Database.Create(&t)
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
