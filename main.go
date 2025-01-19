package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"
)

var BaseDir = os.Getenv("ATHOCS_BASE_DIR")

func main() {
    app := fiber.New()

    // web app
    app.Get("/portal/*", static.New(filepath.Join(BaseDir, "frontend")))

    // graphs
    app.Get("/graphs/*", static.New(filepath.Join(BaseDir, "graphs")))

    // api
    api := app.Group("/api")
    api.Post("/upload", UploadHandler)
    api.Get("/data", DataHandler)

    log.Fatal(app.Listen(":3000"))
}
