# Project Ladybug 

The simple and straightforward testcase management tools.

## Description

TODO introduce features
The Booking sample app demonstrates ([browse the source](https://github.com/revel/samples/tree/master/booking)):

* Using an SQL (SQLite) database and configuring the Revel DB module.
* Using the third party [GORP](https://github.com/coopernurse/gorp) *ORM-ish* library
* [Interceptors](../manual/interceptors.html) for checking that an user is logged in.
* Using [validation](../manual/validation) and displaying inline errors


	booking/app/
		models		   # Structs and validation.
			booking.go
			hotel.go
			user.go

		controllers
			init.go    # Register all of the interceptors.
			gorp.go    # A plugin for setting up Gorp, creating tables, and managing transactions.
			app.go     # "Login" and "Register new user" pages
			hotels.go  # Hotel searching and booking

		views
			...


## I need any kind of help! 

Since I'm new to Go language, not familiar with code convention, documentation, making excellemt code of Go language. Good at HTML/CSS/Javascript? Not at all! I don't have any chance to join web project since I start to work. So the code and poor design may(or shall?) disappoint you. Do not stay. Please make an issue or fork this repo, pull request. Every issues and pull requests are always welcome.

## Getting Started

### Installation

The Project Ladybug uses below... 

* [revel](https://github.com/revel/revel) web framework.
* [gorm](https://github.com/jinzhu/gorm) database driver. 


### Databases

This app uses now only Postgresql. Various databases(MySQL, MarinaDB ....)  will be supported. 

## sqlite Installation

The booking app uses [go-sqlite3](https://github.com/mattn/go-sqlite3) database driver (which wraps the native C library). 


### To install on OSX:

1. Install [Homebrew](http://mxcl.github.com/homebrew/) if you don't already have it.
2. Install pkg-config and sqlite3:

~~~
$ brew install pkgconfig sqlite3
~~~

### To install on Ubuntu:

	$ sudo apt-get install sqlite3 libsqlite3-dev

Once you have SQLite installed, it will be possible to run the booking app:

	$ revel run github.com/revel/samples/booking


## Interceptors

[`app/controllers/init.go`](https://github.com/revel/samples/blob/master/booking/app/controllers/init.go) 
registers the [interceptors](../manual/interceptors.html) that run before every action:

{% highlight go %}
func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
{% endhighlight %}

As an example, `checkUser` looks up the username in the session and redirects
the user to log in if they are not already.

{% highlight go %}
func (c Hotels) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(Application.Index)
	}
	return nil
}
{% endhighlight %}

[Check out the user management code in app.go](https://github.com/revel/samples/blob/master/booking/app/controllers/app.go)

### Features Next

* Requirements management
* Issue management
* Test coverage
