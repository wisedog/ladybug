# Project Ladybug 

The simple and straightforward testcase management tools.

[![Build Status](https://secure.travis-ci.org/wisedog/ladybug.svg?branch=master)](http://travis-ci.org/wisedog/ladybug)
[![Code Climate](https://codeclimate.com/github/wisedog/ladybug/badges/gpa.svg)](https://codeclimate.com/github/wisedog/ladybug)
[![Open Source Love](https://badges.frapsoft.com/os/gpl/gpl.svg?v=102)](https://github.com/ellerbrock/open-source-badge/)

## Description

Project Ladybug can 

* support dashboard
* manage test case
* manage builds
* manage requirements
* support reports(soon)

## I need any kind of help! 

Since I'm new to Go language, not familiar with code convention, documentation, making excellemt code of Go language. Good at HTML/CSS/Javascript? Not at all! I don't have any chance to join web project since I start to work. So the code and poor design may(or shall?) disappoint you. Do not stay. Please make an issue or fork this repo, pull request. Every issues and pull requests are always welcome.

## Getting Started

### Prerequirements

The Project Ladybug uses below... 

* [go](http://www.golang.org) v1.6 or higher
* [gorilla toolkits](http://www.gorillatoolkit.org)
* [gorm](https://github.com/jinzhu/gorm) database driver.
* [bower](http://www.wwwwwwww.org)  

### Installation

* You need set up database before running ladybug
* Now only Postgresql is supported. You can use not only Postgresql but other relational database, but not tested. 
* Default database name is "ladybug" and user "ladybug". see "gorm.go" in app/controllers/gorm.go

### Databases

This app uses now only Postgresql. Various databases(MySQL, MarinaDB ....) will be supported. 

### Run

for development. Try to connect http://localhost:8000 after executing commands below 
```
$ bower install
$ go get -t -d -v ./... 
$ go build
$ ./ladybug 
```

for deployment to Heroku 
```
$ ./ladybug -mode=dev -addr= -port=$PORT -db=$DATABASE_URL
``` 

### Features Next
* Reports
* Test environment
* Milestone
* Test coverage

### Open Sources
* Gorilla toolkits : http://www.gorillatoolkit.org
* gorm : https://github.com/jinzhu/gorm
* AdminLTE : https://github.com/almasaeed2010/AdminLTE
* Icon : http://www.freestockphotos.biz/stockphoto/10655
