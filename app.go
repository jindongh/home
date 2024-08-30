package main

import (
    "log"
    "os"
    "path/filepath"


    "github.com/joho/godotenv"
    "github.com/gofiber/fiber/v3"
)

func main() {
    loadConfig()
    port := os.Getenv("PORT")

    // Initialize a new Fiber app
    app := fiber.New()

    // Define a route for the GET method on the root path '/'
    app.Get("/", func(c fiber.Ctx) error {
        // Send a string response to the client
        return c.SendString("Hello, World 👋!")
    })

    // Start the server on port 
    log.Fatal(app.Listen("0.0.0.0:" + port))
}

func loadConfig() {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }

    // Get the directory of the executable.
    dir := filepath.Dir(ex)
    godotenv.Load(dir + "/.env")
}
