package main

import (
	"bufio"
	"encoding/csv"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func SendCSV[T CSVRowData](c fiber.Ctx, data []T) error {
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="data.csv"`)

	return c.SendStreamWriter(func(w *bufio.Writer) {
		csvWriter := csv.NewWriter(w)
		defer csvWriter.Flush()

		if len(data) == 0 {
			return
		}

		if err := csvWriter.Write(data[0].Header()); err != nil {
			fmt.Fprintf(w, "Error writing CSV header: %v\n", err)
			return
		}

		for _, d := range data {
			if d.Hidden() {
				continue
			}

			if err := csvWriter.Write(d.ToCSVRow()); err != nil {
				fmt.Fprintf(w, "Error writing CSV row: %v\n", err)
				return
			}
		}
	})
}
