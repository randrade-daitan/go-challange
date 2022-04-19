package orm

import (
	"os"
	"testing"
)

func TestOrmDSN(t *testing.T) {
	os.Setenv("DBUSER", "aaa")
	os.Setenv("DBPASS", "bbb")

	dsn := dbDSN()
	expectedDSN := "aaa:bbb@tcp(127.0.0.1:3306)/challange?charset=utf8&parseTime=True&loc=Local"

	if dsn != expectedDSN {
		t.Errorf("orm dsn is incorrect")
	}
}
