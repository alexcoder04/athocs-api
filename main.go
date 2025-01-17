package main

import (
    "log"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
    app := fiber.New()

    // web app
    app.Get("/portal/*", static.New("../athocs-frontend/build"))

    // api
    api := app.Group("/api")
    api.Post("/upload", UploadHandler)
    api.Get("/data", DataHandler)

    log.Fatal(app.Listen(":3000"))
}
