package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type EmptyRowError struct{}

func (e *EmptyRowError) Error() string {
	return "This is an empty row"
}

func ReadCSV[T CSVRowData](p string, parse func([]string) (T, error)) ([]T, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(f)

	// read header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var data []T

	for {
		record, err := reader.Read()
		if err == io.EOF { // we are done
			break
		}
		if err != nil { // failed to read row
			return data, err
		}

		// parse data using given function
		d, err := parse(record)
		if err != nil {
			if errors.Is(err, &EmptyRowError{}) {
				continue
			}
			return data, err
		}
		data = append(data, d)
	}

	return data, nil
}

func WriteToCSV(p string, row []string) error {
	f, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	return writer.Write(row)
}
