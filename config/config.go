package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var configurations *Config

type Config struct {
	Version       string
	ServiceName   string
	HttpPort      int
	JwtSecretKey  string
	JwtRefreshKey string
	DbHost        string
	DbPort        int
	DbUser        string
	DbPassword    string
	DbName        string
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load env variables: ", err)
		os.Exit(1)
	}
	version := os.Getenv("VERSION")
	if version == "" {
		fmt.Println("Version is required!")
		os.Exit(1)
	}
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		fmt.Println("Service name is required!")
		os.Exit(1)
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		fmt.Println("Http Port is required!")
		os.Exit(1)
	}
	port, err := strconv.ParseInt(httpPort, 10, 64)
	if err != nil {
		fmt.Println("Port must be number")
		os.Exit(1)
	}
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		fmt.Println("Jwt secret key is required!")
		os.Exit(1)
	}
	jwtRefreshKey := os.Getenv("JWT_REFRESH_KEY")
	if jwtRefreshKey == "" {
		fmt.Println("Jwt refresh key is required!")
		os.Exit(1)
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	db_port, err := strconv.ParseInt(dbPort, 10, 64)
	if err != nil {
		fmt.Println("DB Port must be a number")
		os.Exit(1)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	configurations = &Config{
		Version:       version,
		ServiceName:   serviceName,
		HttpPort:      int(port),
		JwtSecretKey:  jwtSecretKey,
		JwtRefreshKey: jwtRefreshKey,
		DbHost:        dbHost,
		DbPort:        int(db_port),
		DbUser:        dbUser,
		DbPassword:    dbPassword,
		DbName:        dbName,
	}
}

func GetConfig() *Config {
	if configurations == nil {
		loadConfig()
	}
	return configurations
}
