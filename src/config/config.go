package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// it sucks that golang does it this way but i couldn't find a better one
const WORKING_DIR = "FILL_IN"

// Config func to get env value
func Config(key string) string {
	// If in a test case, use different config file
	if flag.Lookup("test.v") != nil {

		// Get the absolute path of the currently executing file
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}

		// Get the directory of the executable file
		exPath := filepath.Dir(ex)
		fmt.Print(exPath)

		err = godotenv.Load(WORKING_DIR + "\\.test.env")
		if err != nil {
			fmt.Print("Error loading .env file")
			fmt.Print(err)
		}
		return os.Getenv(key)
	}

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
