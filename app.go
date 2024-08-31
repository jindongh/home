package main

import (
    "embed"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/markbates/goth"
    "github.com/jindongh/home/docker"
    "github.com/markbates/goth/providers/google"
    "github.com/shareed2k/goth_fiber"
    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var viewsfs embed.FS
//go:embed static/*
var staticfs embed.FS

type config struct {
    port string
    clientId string
    clientSecret string
    HomeUrl string
    Services []docker.ServiceConfig
}
func main() {
    config := loadConfig()
    store := session.New()
    dockerService := docker.NewService()

    // Initialize a new Fiber app
    engine := html.NewFileSystem(http.FS(viewsfs), ".html")
    engine.AddFuncMap(map[string]interface{}{
        "Title": func(s string) string {
            return strings.Title(s)
        },
    })
    app := fiber.New(fiber.Config{
        Views: engine,
    })
    app.Use("/static", filesystem.New(filesystem.Config{
        Root: http.FS(staticfs),
        PathPrefix: "static",
    }))
    app.Get("/pi/:action?", func(c *fiber.Ctx) error {
        sess, _ := store.Get(c)
        return c.Render("templates/pi/" + c.Params("action"), fiber.Map{
            "Email": sess.Get("email"),
            "Config": config,
        })
    })
    app.Get("/admin", func(c *fiber.Ctx) error {
        sess, _ := store.Get(c)
        return c.Render("templates/admin", fiber.Map{
            "Email": sess.Get("email"),
            "Config": config,
            "Services": dockerService.GetServiceStatus(),
        })
    })
    app.Get("/admin/:service?/:action?", func(ctx *fiber.Ctx) error {
        service := ctx.Params("service")
        if ctx.Params("action") == "up" {
            dockerService.StartService(service)
        } else {
            dockerService.StopService(service)
        }
        return ctx.Redirect("/admin")
    })
    app.Get("/", func(c *fiber.Ctx) error {
        sess, _ := store.Get(c)
        return c.Render("templates/index", fiber.Map{
            "Email": sess.Get("email"),
            "Config": config,
        })
    })

    goth.UseProviders(
        google.New(config.clientId, config.clientSecret, config.HomeUrl + "/auth/callback"),
    )
    app.Get("/login", goth_fiber.BeginAuthHandler)
    app.Get("/auth/callback", func(ctx *fiber.Ctx) error {
        user, err := goth_fiber.CompleteUserAuth(ctx)
        if err != nil {
            log.Fatal(err)
        }
        sess, _ := store.Get(ctx)
        sess.Set("email", user.Email)
        sess.Save()
        return ctx.Redirect("/")
    })
    app.Get("/logout", func(ctx *fiber.Ctx) error {
        if err := goth_fiber.Logout(ctx); err != nil {
            log.Fatal(err)
        }
        sess, _ := store.Get(ctx)
        sess.Destroy()
        return ctx.Redirect("/")
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
        clientId: os.Getenv("GOOGLE_CLIENT_ID"),
        clientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        HomeUrl: os.Getenv("URL_HOME"),
        Services: docker.GetServiceConfigs(),
    }
}

