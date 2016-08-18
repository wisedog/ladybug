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
