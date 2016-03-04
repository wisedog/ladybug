package controllers

import (
	"github.com/wisedog/ladybug/models"
	"github.com/wisedog/ladybug/app/routes"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type Application struct {
	GormController
}

func (c Application) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.RenderArgs["user"] = user
	}
	return nil
}

func (c Application) connected() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c Application) getUser(username string) *models.User {
	/*users, err := c.Tx.Select(models.User{}, `select * from User where Email = ?`, username)

	if err != nil {
		panic(err)
	}*/
	user := models.User{}
	c.Tx.Where("Email = ?", username).First(&user)
	/*if len(user) == 0 {
		return nil
	}*/
	return &user
}

func (c Application) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Hello.Welcome())
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c Application) Register() revel.Result {
	return c.Render()
}

func (c Application) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Application.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	c.Tx.NewRecord(user)
	c.Tx.Create(user)

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.Users.Index(user.ID))
}

func (c Application) Login(email, password string, remember bool) revel.Result {
	user := c.getUser(email)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = email
			c.Session["user_id"] = strconv.Itoa(user.ID)
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + email)

			// Set language from user's preference.
			// TODO add language field to User model
			// language field should be string and ISO639-1 codes.
			// Region field should be string and ISO3166-1 alpha-2 code
			
			cookie := http.Cookie{Name: "REVEL_LANG", Value: "en-US", Path: "/"}
			c.SetCookie(&cookie)
			
			user.LastLoginAt = time.Now()
			c.Tx.Save(&user)
			return c.Redirect(routes.Hello.Welcome())
		}
	}

	c.Flash.Out["user"] = email
	c.Flash.Error("Login failed")
	return c.Redirect(routes.Application.Index())
}

func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Index())
}

// Return from status code to string. 
// TODO should be localization
func (c Application) getStatusL10N(status int) string{
	str := "Unknown Status"
	switch{
		case status == models.PRIORITY_HIGHEST:
			str = "Highest"
		case status == models.PRIORITY_HIGH:
			str = "High"
		case status == models.PRIORITY_MEDIUM:
			str = "Medium"
		case status == models.PRIORITY_LOW:
			str = "Low"
		case status == models.PRIORITY_LOWEST:
			str = "Lowest"
	}
	
	
	return str
}