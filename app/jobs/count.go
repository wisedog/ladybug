package jobs

import (
	"fmt"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

// Periodically count the bookings in the database.
type BookingCounter struct{}

func (c BookingCounter) Run() {
	//	bookings := []models.Booking{}
	//	controllers.Db.Find(&bookings)
	/*	bookings, err := controllers.Dbm.Select(models.Booking{},
			`select * from Booking`)
		if err != nil {
			panic(err)
		}*/
	fmt.Printf("Heartbeat \n")
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 10m", BookingCounter{})
	})
}
