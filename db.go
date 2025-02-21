package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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

func WriteDatapoint(data *Datapoint) error {
	filename, err := GetCurrentFile(data.Timestamp)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := InitNewFile(filename)
		if err != nil {
			return err
		}
	}

	err = WriteToCSV(filename, data.ToCSVRow())
	if err != nil {
		return err
	}

	return AddStation(data.Station)
}

func ReadDataForStation(station string, date string, start time.Time, end time.Time) ([]Datapoint, error) {
	// we pass a function that tells how to parse the csv items
	data, err := ReadCSV(filepath.Join(Config.DBDir, date+".csv"), func(row []string) (Datapoint, error) {
		// parse data
		dp, err := DatapointFromCSV(row)
		if err != nil {
			return dp, err
		}

		// check whether we want that datapoint
		timestamp, _ := time.Parse(Config.TimestampFormat, dp.Timestamp)
		if dp.Station == station && timestamp.After(start) && timestamp.Before(end) {
			return dp, nil
		}
		// returning an EmptyRowError will make the csv parser discard that datapoint
		return dp, &EmptyRowError{}
	})

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

func FetchData(req *DataRequest) ([]Datapoint, error) {
	var end time.Time
	var err error
	if req.TimeTo == "" {
		end = time.Now()
	} else {
		end, err = time.Parse(Config.TimestampFormat, req.TimeTo)
		if err != nil {
			return nil, err
		}
	}

	var start time.Time
	if req.TimeFrom == "" {
		start = end.Add(-time.Duration(Config.DefaultDataInterval) * time.Hour)
	} else {
		start, err = time.Parse(Config.TimestampFormat, req.TimeFrom)
		if err != nil {
			return nil, err
		}
	}

	if start.After(end) {
		return nil, fmt.Errorf("start timestamp must be before or equal to end timestamp")
	}

	data := []Datapoint{}
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		curData, err := ReadDataForStation(req.Station, current.Format(Config.DateFormat), start, end)
		if err != nil {
			return nil, err
		}

		data = append(data, curData...)
	}

	return data, nil
}
