package utils

import "github.com/joho/godotenv"

func LoadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		panic("Error loading " + path + " file")
	}
}
