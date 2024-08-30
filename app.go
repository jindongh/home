package main

import (
    "context"
    "embed"
    "log"
    "net/http"
    "os"
    "path/filepath"

    containertypes "github.com/docker/docker/api/types/container"
    apitypes "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/markbates/goth"
    "github.com/markbates/goth/providers/google"
    "github.com/shareed2k/goth_fiber"
    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/gofiber/template/html/v2"
)
var downloadService = []string{"/aria2-pro", "/ariang"}

//go:embed templates/*
var viewsfs embed.FS
//go:embed static/*
var staticfs embed.FS

type config struct {
    port string
    clientId string
    clientSecret string
    HomeUrl string
    BookUrl string
    DownloadUrl string
    VideoUrl string
    PhotoUrl string
}
type service struct {
    IsDownloadUp bool
}
func main() {
    config := loadConfig()
    store := session.New()

    // Initialize a new Fiber app
    app := fiber.New(fiber.Config{
        Views: html.NewFileSystem(http.FS(viewsfs), ".html"),
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
            "Service": getServices(),
        })
    })
    app.Get("/admin/:service?/:action?", func(ctx *fiber.Ctx) error {
        if ctx.Params("action") == "up" {
            startService(downloadService)
        } else {
            stopService(downloadService)
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
        VideoUrl: os.Getenv("URL_VIDEO"),
        PhotoUrl: os.Getenv("URL_PHOTO"),
        BookUrl: os.Getenv("URL_VIDEO"),
        DownloadUrl: os.Getenv("URL_DOWNLOAD"),
    }
}
func getServices() *service {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        panic(err)
    }
    defer cli.Close()

    if err != nil {
        panic(err)
    }
    containers := findService(cli, downloadService)
    isUp := true
    for _, container := range(containers) {
        if container.State != "running" {
            log.Println("found container not running", container.Names, container.State)
            isUp = false
        }
    }
    return &service{
        IsDownloadUp: isUp,
    }
}
func startService(names []string) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        panic(err)
    }
    defer cli.Close()
    containers := findService(cli, names)
    for _, container := range(containers) {
        log.Println("begin start container", container.Names)
        err = cli.ContainerStart(ctx, container.ID, containertypes.StartOptions{})
        log.Println("finish start container", container.Names, err)
    }
}
func stopService(names []string) {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        panic(err)
    }
    defer cli.Close()
    containers := findService(cli, names)
    log.Println("begin stop containers", len(containers))
    for _, container := range(containers) {
        log.Println("begin stop container", container.Names)
        err = cli.ContainerStop(ctx, container.ID, containertypes.StopOptions{})
        log.Println("finish stop container", container.Names, err)
    }
}

func findService(cli *client.Client, names []string) []*apitypes.Container {
    ctx := context.Background()
    containers, err := cli.ContainerList(ctx, containertypes.ListOptions{All: true})
    if err != nil {
        log.Fatalf("failed to find container", err)
    }
    matchContainers := []*apitypes.Container{}
    for _, container := range containers {
        for _, containerName := range(container.Names) {
            for _, name := range(names) {
                if containerName == name {
                    log.Println("found match container", name, container.ID)
                    matchContainers = append(matchContainers, &container)
                }
            }
        }
    }
    return matchContainers
}
