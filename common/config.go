package common

import (
	"os"

	"github.com/joho/godotenv"
	"path/filepath"
	"github.com/jindongh/home/docker"
)

type Config struct {
    Port string
    ClientId string
    ClientSecret string
    HomeUrl string
    Services []docker.ServiceConfig
}
func LoadConfig() *Config {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }

    // Get the directory of the executable.
    dir := filepath.Dir(ex)
    godotenv.Load(dir + "/.env")
    return &Config{
        Port: os.Getenv("PORT"),
        ClientId: os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        HomeUrl: os.Getenv("URL_HOME"),
        Services: docker.GetServiceConfigs(),
    }
}

