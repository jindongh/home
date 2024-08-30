package main

import (
    "embed"
    "log"
    "net/http"
    "os"
    "path/filepath"


    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    "github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var viewsfs embed.FS
//go:embed static/*
var staticfs embed.FS

type config struct {
    port string
    HomeUrl string
    BookUrl string
    DownloadUrl string
    VideoUrl string
    PhotoUrl string
}
func main() {
    config := loadConfig()

    // Initialize a new Fiber app
    app := fiber.New(fiber.Config{
        Views: html.NewFileSystem(http.FS(viewsfs), ".html"),
    })
    app.Use("/static", filesystem.New(filesystem.Config{
        Root: http.FS(staticfs),
        PathPrefix: "static",
    }))
    app.Get("/pi/:action?", func(c *fiber.Ctx) error {
        return c.Render("templates/pi/" + c.Params("action"), fiber.Map{
            "Config": config,
        })
    })
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("templates/index", fiber.Map{
            "Config": config,
        })
    })

    // Start the server on port 
    log.Fatal(app.Listen("0.0.0.0:" + config.port))
}

func loadConfig() *config {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }

    // Get the directory of the executable.
    dir := filepath.Dir(ex)
    godotenv.Load(dir + "/.env")
    return &config{
        port: os.Getenv("PORT"),
        HomeUrl: os.Getenv("URL_HOME"),
        VideoUrl: os.Getenv("URL_VIDEO"),
        PhotoUrl: os.Getenv("URL_PHOTO"),
        BookUrl: os.Getenv("URL_VIDEO"),
        DownloadUrl: os.Getenv("URL_DOWNLOAD"),
    }
}
