package main

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

// get list of existing stations
func StationsListHandler(c fiber.Ctx) error {
	stations, err := GetStationsList()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load stations list from database",
		})
	}

	return SendCSV(c, stations)
}

// fetch data from db
func DataHandler(c fiber.Ctx) error {
	req := new(DataRequest)

	if err := c.Bind().Query(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	data, err := FetchData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to load data from database",
			"goError": err.Error(),
		})
	}

	return SendCSV(c, data)
}

// upload new data to db
func UploadHandler(c fiber.Ctx) error {
	data := new(Datapoint)

	if err := c.Bind().Body(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data",
		})
	}

	// esp32 does not have rtc
	if data.Timestamp == "auto" {
		data.Timestamp = time.Now().Format("2006-01-02_15:04:05")
	}

	if err := WriteDatapoint(data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to write data to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data successfully saved",
	})
}
