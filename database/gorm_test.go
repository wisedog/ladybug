package database

import(
	"testing"

  "github.com/wisedog/ladybug/interfacer"
)

func TestDialect(t *testing.T) {

  conf := interfacer.LoadConfig()
  rv := getDialectArgs(conf)
  
  if rv != "postgres://ladybug:@localhost:5432/ladybug?sslmode=false"{
    t.Error("failed")
  }
}
