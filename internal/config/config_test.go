package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	res, err := LoadConfig()
	if err != nil {
		t.Errorf("There was an error loading")
	}

	if res.Port != PORT {
		t.Errorf("Port is incorrect")
	}
	if res.DB_Url != DB_URL {
		t.Errorf("DBURL is incorrect")
	}
}
