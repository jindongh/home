package main

import (
    "embed"
    "log"
    "net/http"
    "strings"

    "github.com/markbates/goth"
    "github.com/jindongh/home/common"
    "github.com/jindongh/home/docker"
    "github.com/jindongh/home/piano"
    "github.com/markbates/goth/providers/google"
    "github.com/shareed2k/goth_fiber"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var viewsfs embed.FS
//go:embed static/*
var staticfs embed.FS

func main() {
    config := common.LoadConfig()
    piano.Connect()
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
    app.Get("/openwrt/:action?", func(c *fiber.Ctx) error {
        sess, _ := store.Get(c)
        return c.Render("templates/openwrt/" + c.Params("action"), fiber.Map{
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
    piano.Route(app, store, config)
    goth.UseProviders(
        google.New(config.ClientId, config.ClientSecret, config.HomeUrl + "/auth/callback"),
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
    log.Fatal(app.Listen("0.0.0.0:" + config.Port))
}
