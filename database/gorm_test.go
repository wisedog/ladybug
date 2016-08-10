package database

import (
	"testing"

	"github.com/wisedog/ladybug/interfacer"
)

// TestDialect test each DB's dialect.
// But now there is only postgresql.
func TestDialect(t *testing.T) {

	conf := interfacer.LoadConfig("test", "./dummy_ladybug.conf")
	rv := getDialectArgs(conf)

	expected := "postgres://ladybug:@localhost:5432/ladybug?sslmode=disable"
	if rv != expected {
		t.Error("expected : ", expected, ", but ", rv)
	}
}

// TestDialectWithArgs test arguments passing db's dialect
func TestDialectWithArgs(t *testing.T) {

	conf := interfacer.LoadConfigWithArgs("test", "localhost", 8000, "testtest", "./dummy_ladybug.conf")
	rv := getDialectArgs(conf)

	expected := "testtest"
	if rv != expected {
		t.Error("expected : ", expected, ", but ", rv)
	}
}
