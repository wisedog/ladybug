package controllers

import (
	//	"fmt"
	"github.com/revel/revel"
	"github.com/wisedog/ladybug/app/models"
	"strings"
	//	"github.com/wisedog/ladybug/app/routes"
	/*		"golang.org/x/crypto/bcrypt"
			"strings"*/)

type Issues struct {
	Application
}

func (c Issues) Detail(project string, id int) revel.Result {
	var issues []models.Issue

	revel.INFO.Println("AAA", project, id)
	return c.Render(issues)

}

func (c Issues) List(project string, search string, size, page int) revel.Result {
	if page == 0 {
		page = 1
	}
	revel.INFO.Println("BBB", project)
	nextPage := page + 1
	search = strings.TrimSpace(search)

	if size == 0 {
		size = 20
	}

	var is []models.Issue
	if search == "" {
		c.Tx.Preload("Assignee").Limit(size).Offset((page - 1) * size).Find(&is)
	} else {
		//only latin-based characters only. search = strings.ToLower(search)
		c.Tx.Preload("Assignee").Limit(size).Offset((page-1)*size).
			Where("summary like ? or description like ?", "%"+search+"%", "%"+search+"%").
			Find(&is)
	}

	return c.Render(is, search, size, page, nextPage)
}

/*
type Hotels struct {
	Application
}

func (c Hotels) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Hotels) Index() revel.Result {
	var bookings []models.Booking

	c.Tx.Preload("User").Preload("Hotel").Where("user_id = ?", c.connected().ID).Find(&bookings)

	for _, r := range bookings {
		fmt.Printf("%+v", r)
	}

	return c.Render(bookings)
}

func (c Hotels) List(search string, size, page int) revel.Result {
	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	search = strings.TrimSpace(search)

	var hotels []models.Hotel
	if search == "" {
		c.Tx.Limit(size).Offset((page - 1) * size).Find(&hotels)
	} else {
		search = strings.ToLower(search)
		c.Tx.Limit(size).Offset((page-1)*size).
			Where("lower(Name) like ? or lower(City) like ?", "%"+search+"%", "%"+search+"%").
			Find(&hotels)
	}

	return c.Render(hotels, search, size, page, nextPage)
}

func loadHotels(results []interface{}, err error) []*models.Hotel {
	if err != nil {
		panic(err)
	}
	revel.INFO.Println("Hotel-loadHotels ")
	var hotels []*models.Hotel
	for _, r := range results {
		hotels = append(hotels, r.(*models.Hotel))
	}
	return hotels
}

func (c Hotels) loadHotelById(id int) *models.Hotel {
	result := new(models.Hotel)
	c.Tx.Where("id = ?", id).First(&result)
	return result
}

func (c Hotels) Show(id int) revel.Result {
	hotel := c.loadHotelById(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}
	title := hotel.Name
	return c.Render(title, hotel)
}

func (c Hotels) Settings() revel.Result {
	return c.Render()
}

func (c Hotels) SaveSettings(password, verifyPassword string) revel.Result {
	models.ValidatePassword(c.Validation, password)
	c.Validation.Required(verifyPassword).
		Message("Please verify your password")
	c.Validation.Required(verifyPassword == password).
		Message("Your password doesn't match")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(routes.Hotels.Settings())
	}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//db.Model(&user).Where("active = ?", true).Update("name", "hello")
	user := models.User{}
	revel.INFO.Println("Hotel-SaveSettings ")
	c.Tx.Model(&user).Where("id = ?", c.connected().ID).Update("HashedPassword", bcryptPassword)
	c.Flash.Success("Password updated")
	return c.Redirect(routes.Hotels.Index())
}

func (c Hotels) ConfirmBooking(id int, booking models.Booking) revel.Result {
	hotel := c.loadHotelById(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}

	title := fmt.Sprintf("Confirm %s booking", hotel.Name)
	booking.Hotel = *hotel
	booking.User = *(c.connected())
	booking.Validate(c.Validation)

	if c.Validation.HasErrors() || c.Params.Get("revise") != "" {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Hotels.Book(id))
	}

	if c.Params.Get("confirm") != "" {
		c.Tx.NewRecord(booking)
		c.Tx.Create(&booking)
		c.Flash.Success("Thank you, %s, your confirmation number for %s is %d",
			booking.User.Name, hotel.Name, booking.ID)
		return c.Redirect(routes.Hotels.Index())
	}

	return c.Render(title, hotel, booking)
}

func (c Hotels) CancelBooking(id int) revel.Result {
	c.Tx.Where("id = ?", id).Delete(models.Booking{})
	c.Flash.Success(fmt.Sprintln("Booking cancelled for confirmation number", id))
	return c.Redirect(routes.Hotels.Index())
}

func (c Hotels) Book(id int) revel.Result {
	hotel := c.loadHotelById(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}

	title := "Book " + hotel.Name
	return c.Render(title, hotel)
}*/
