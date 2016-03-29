package interfacer

import(
	"testing"
)

func TestLoadConfig(t *testing.T) {
  cf := LoadConfig()

  if cf == nil{
    t.Error("config load failed")
  }
}


func TestMode(t *testing.T) {
  cf := LoadConfig()

  if cf == nil{
    t.Error("config load failed")
  }

  if cf.GetMode() != "dev"{
    t.Error("Config value is not property")
  }
}

func TestBindAddress(t *testing.T){
  cf := LoadConfig()

  if cf == nil{
    t.Error("config load failed")
  }

  if cf.GetGetBindAddress() != "localhost:8000"{
    t.Error("Config value is not property")
  }
}