package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionString(t *testing.T) {
	c := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "password",
		Database: "dbname",
		SSLMode:  "disable",
	}
	expected := "host=localhost port=5432 user=user password=password dbname=dbname sslmode=disable TimeZone=UTC"
	actual := ConnectionDBString(c)
	assert.Equal(t, expected, actual, "Invalid format")
}
