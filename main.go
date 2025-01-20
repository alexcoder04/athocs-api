package main

import (
    "log"
    "path/filepath"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
    app := fiber.New()

    // web app
    app.Get("/portal/*", static.New(filepath.Join(Config.BaseDir, "frontend")))

    // graphs
    app.Get("/graphs/*", static.New(filepath.Join(Config.BaseDir, "graphs")))

    // api
    api := app.Group("/api")
    api.Post("/upload", UploadHandler)
    api.Get("/data", DataHandler)

    log.Fatal(app.Listen(":"+Config.Port))
}
