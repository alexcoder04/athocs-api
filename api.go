package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"strconv"
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

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="data.csv"`)

	return c.SendStreamWriter(func(w *bufio.Writer) {
		csvWriter := csv.NewWriter(w)

		header := []string{"id"}
		if err := csvWriter.Write(header); err != nil {
			fmt.Fprintf(w, "Error writing CSV header: %v\n", err)
			return
		}

		for _, st := range stations {
			row := []string{st}
			if err := csvWriter.Write(row); err != nil {
				fmt.Fprintf(w, "Error writing CSV row: %v\n", err)
				return
			}
		}

		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			return
		}
	})
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

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="data.csv"`)

	return c.SendStreamWriter(func(w *bufio.Writer) {
		csvWriter := csv.NewWriter(w)

		header := []string{"timestamp", "station", "temperature", "humidity", "pressure", "battery"}
		if err := csvWriter.Write(header); err != nil {
			fmt.Fprintf(w, "Error writing CSV header: %v\n", err)
			return
		}

		for _, dp := range data {
			row := []string{
				dp.Timestamp,
				dp.Station,
				strconv.FormatFloat(float64(dp.Temperature), 'f', 2, 32),
				strconv.FormatFloat(float64(dp.Humidity), 'f', 2, 32),
				strconv.FormatFloat(float64(dp.Pressure), 'f', 2, 32),
				strconv.Itoa(int(dp.Battery)),
			}
			if err := csvWriter.Write(row); err != nil {
				fmt.Fprintf(w, "Error writing CSV row: %v\n", err)
				return
			}
		}

		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			return
		}
	})
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
