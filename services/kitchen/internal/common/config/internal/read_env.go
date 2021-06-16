package internal

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// ReadEnv into variable
func ReadEnv(cfg interface{}) {
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatal("Error reading env variables", "\n\t", err)
	}
}
