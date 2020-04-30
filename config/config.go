package config

import (
	"log"
	"os"
)

var reader yamlReader

func init() {
	log.Println("Initializing config...")
	err := reader.load()
	if err != nil {
		log.Println("Failed to load application.yaml, error: ", err)
		os.Exit(1)
	}
	reader.exec()
}

func Exec() {
	reader.exec()
}

func GetString(name string) string {
	return reader.getString(name)
}

func GetInt(name string) int {
	return reader.getInt(name)
}

func GetBool(name string) bool {
	return reader.getBool(name)
}

func GetFloat64(name string) float64 {
	return reader.getFloat64(name)
}
