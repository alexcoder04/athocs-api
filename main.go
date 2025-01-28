package main

import (
    "log"
    "path/filepath"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
    app := fiber.New()

    // api
    app.Get("/api/data", DataHandler)
    app.Post("/api/upload", UploadHandler)

    // graphs
    app.Get("/graphs/*", static.New(filepath.Join(Config.BaseDir, "graphs")))

    // web app
    app.Get("/*", static.New(filepath.Join(Config.BaseDir, "frontend")))

    log.Fatal(app.Listen(":"+Config.Port))
}
