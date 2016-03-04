package controllers

import (
  "fmt"
  "testing"
  "os"
  )

func mySetupFunction(){
  fmt.Println("Setup...")
}

func myTeardownFunction(){
  fmt.Println("Tear Down...")
}

func TestMain(m *testing.M) { 
  fmt.Println("TestMain...")
  mySetupFunction()
  retCode := m.Run()
  myTeardownFunction()
  os.Exit(retCode)
}

