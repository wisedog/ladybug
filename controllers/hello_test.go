package controllers

import (
	"fmt"
	"os"
	"testing"
)

func mySetupFunction() {
	fmt.Println("Setup...")
}

func myTeardownFunction() {
	fmt.Println("Tear Down...")
}

func TestMain(m *testing.M) {
	fmt.Println("TestMain...")
	mySetupFunction()
	retCode := m.Run()
	myTeardownFunction()
	os.Exit(retCode)
}

func TestTemp1(t *testing.T) {}

func TestTemp2(t *testing.T) {}
func TestTemp3(t *testing.T) {}
func TestTemp4(t *testing.T) {}
func TestTemp5(t *testing.T) {}
func TestTemp6(t *testing.T) {}
