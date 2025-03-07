package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	debug := flag.Bool("debug", false, "enable cors for all origins")
	flag.Parse()

	app := fiber.New()

	if *debug {
		// allowing all origins is "bad practice" and disallowed in fiber
		// we dont care and use AllowOriginsFunc to bypass that restriction
		app.Use(cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				return true
			},
		}))
	}

	err := LiveDB.Init()
	if err != nil {
		log.Fatalf("Failed to init LiveDB, %s\n", err.Error())
	}

	// api
	app.Get("/api/stations", StationsListHandler)
	app.Get("/api/data", DataHandler)
	app.Get("/api/data/live", DataLiveHandler)
	app.Post("/api/upload", UploadHandler)

	// graphs
	app.Get("/graphs/*", static.New(filepath.Join(Config.BaseDir, "graphs")))

	// web app
	app.Get("/*", static.New(filepath.Join(Config.BaseDir, "frontend")))

	log.Fatal(app.Listen(":" + Config.Port))
}
