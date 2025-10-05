package utils

import (
	"encoding/json"
	"log"
)

func ContainsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

var LogMode = "dev" // "dev" ou "prod"

func Logger(msg string, data ...interface{}) {
	if len(data) == 0 {
		log.Println(msg)
		return
	}

	for _, d := range data {
		pretty, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			log.Printf("%s: %v\n", msg, d)
		} else {
			log.Printf("%s: %s\n", msg, string(pretty))
		}
	}
}
