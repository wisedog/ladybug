package database

import(
	"testing"

  "github.com/wisedog/ladybug/interfacer"
)

func TestDialect(t *testing.T) {

  conf := interfacer.LoadConfig()
  rv := getDialectArgs(conf)

  expected := "postgres://ladybug:@localhost:5432/ladybug?sslmode=disable"
  if rv != expected{
    t.Error("expected : ", expected, ", but ", rv)
  }
}
