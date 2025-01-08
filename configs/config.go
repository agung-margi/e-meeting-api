package configs

import (
	"encoding/json"
	"os"
)

var AppConfig Config

func LoadConfig() {
	file, err := os.Open("configs/config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		panic(err)
	}
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Timezone string `json:"timezone"`
}

type Config struct {
	Database     DatabaseConfig `json:"database"`
	JWTSecretKey string         `json:"jwt_secret_key"`
	BaseURL      string         `json:"base_url"`
}
