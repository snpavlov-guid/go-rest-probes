package gormdb

import (
	"testing"

)

func TestString(t *testing.T) {

	gctx := GormDBContext{}

	dsn := "host=localhost port=5432 dbname=GormTestDB user=RIMDBAdmin password=RimDBAdmin123 sslmode=disable"

	err := gctx.Migrate(dsn, "lpok")

	if err != nil {
		t.Fatalf("Cant migrate database! Error: %v", err)
	}

}
