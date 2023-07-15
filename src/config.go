package src

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Loaded environment variables from .env file")
	Envars = Envs{
		Token:     os.Getenv("TOKEN"),
		DbUrl:     os.Getenv("DB_URL"),
		Encyptkey: os.Getenv("ENCRYPT_KEY"),
		Port:      os.Getenv("PORT"),
	}

}

var Envars = Envs{}

func Converttoin32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int32(i)
}

func ConverttoFloat64(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return i
}
