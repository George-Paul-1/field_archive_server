package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	// IN BETWEEN THESE LINES IS A FIX FOR THE TEST EXECUTING IN A DIFFERENT MASTER DIRECTORY TO THE PROGRAM.
	// -------------------------------------------------------
	// Change the working directory to the root directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	// Change to root directory where .env file is located
	err = os.Chdir("../../") // Assuming root directory is one level above
	if err != nil {
		t.Fatalf("Failed to change directory to root: %v", err)
	}
	defer os.Chdir(originalDir) // Ensure we return to the original directory after the test
	// ---------------------------------------------------------

	res, err := LoadConfig()

	if err != nil {
		t.Errorf("There was an error loading Config %v", err)
	}

	if res.Port != PORT {
		t.Errorf("Port is incorrect")
	}
	if res.DB_Url != DB_URL {
		t.Errorf("DBURL is incorrect")
	}
}
