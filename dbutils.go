package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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

// write csv header
func InitNewFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("timestamp,station,temperature,humidity,pressure,battery\n")
	return err
}

// extract date from full timestamp
func GetCurrentFile(timestamp string) (string, error) {
	parsedDate, err := time.Parse(Config.TimestampFormat, timestamp)
	if err != nil {
		return "", err
	}

	dateString := parsedDate.Format(Config.DateFormat)
	return filepath.Join(Config.DBDir, dateString+".csv"), nil
}

func GetStartEndTime(start string, end string) (time.Time, time.Time, error) {
	var endTime time.Time
	var startTime time.Time
	var err error

	if end == "" {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse(Config.TimestampFormat, end)
		if err != nil {
			return startTime, endTime, err
		}
	}

	if end == "" {
		startTime = time.Now().Add(-time.Duration(Config.DefaultDataInterval) * time.Hour)
	} else {
		startTime, err = time.Parse(Config.TimestampFormat, end)
		if err != nil {
			return startTime, endTime, err
		}
	}

	if startTime.After(endTime) {
		return startTime, endTime, fmt.Errorf("start timestamp must be before end timestamp")
	}

	return startTime, endTime, nil
}
