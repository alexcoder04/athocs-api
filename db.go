package main

import (
	"os"
	"path/filepath"
	"time"
)

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

	// write to persistent db
	err = WriteToCSV(filename, data.ToCSVRow())
	if err != nil {
		return err
	}

	// write to live db
	LiveDB.Set(*data)

	return AddStation(data.Station)
}

func ReadDataForStation(station string, date string, start *time.Time, end *time.Time) ([]Datapoint, error) {
	// we pass a function that tells how to parse the csv items
	data, err := ReadCSV(filepath.Join(Config.DBDir, date+".csv"), func(row []string) (*Datapoint, error) {
		// parse data
		dp, err := DatapointFromCSV(row)
		if err != nil {
			return nil, err
		}

		// check whether we want that datapoint
		timestamp, _ := time.Parse(Config.TimestampFormat, dp.Timestamp)
		if dp.Station == station && timestamp.After(*start) && timestamp.Before(*end) {
			return dp, nil
		}
		return nil, nil
	})

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

// this is the actual function called by the route handler
func FetchData(req *DataRequest) ([]Datapoint, error) {
	start, end, err := GetStartEndTime(req.TimeFrom, req.TimeTo)
	if err != nil {
		return nil, err
	}

	data := []Datapoint{}
	// run through dates
	for current := *start; !current.After(*end); current = current.Add(24 * time.Hour) {
		curData, err := ReadDataForStation(req.Station, current.Format(Config.DateFormat), start, end)
		if err != nil {
			return nil, err
		}

		data = append(data, curData...)
	}

	return data, nil
}

func FetchLatestData() []Datapoint {
	data := []Datapoint{}
	for _, dp := range LiveDB.Data {
		data = append(data, dp)
	}
	return data
}
