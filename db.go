package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alexcoder04/friendly/v2"
)

const (
	TIMESTAMP_FORMAT = "2006-01-02_15:04:05"
	DATE_FORMAT      = "2006-01-02"
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
	parsedDate, err := time.Parse(TIMESTAMP_FORMAT, timestamp)
	if err != nil {
		return "", err
	}

	dateString := parsedDate.Format(DATE_FORMAT)
	return filepath.Join(Config.DBDir, dateString+".csv"), nil
}

func GetStationsList() ([]string, error) {
	f, err := os.Open(filepath.Join(Config.DBDir, "stations-index.csv"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// read header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var stations []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return stations, err
		}

		stations = append(stations, record[0])
	}

	return stations, nil
}

func AddStation(station string) error {
	stations, err := GetStationsList()
	if err != nil {
		return err
	}

	if friendly.ArrayContains(stations, station) {
		return nil
	}

	f, err := os.OpenFile(filepath.Join(Config.DBDir, "stations-index.csv"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf(
		"%s\n",
		station,
	))
	return err
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

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf(
		"%s,%s,%.2f,%.2f,%.1f,%d\n",
		data.Timestamp,
		data.Station,
		data.Temperature,
		data.Humidity,
		data.Pressure,
		data.Battery,
	))
	if err != nil {
		return err
	}

	return AddStation(data.Station)
}

func ReadDataForStation(station string, date string, start time.Time, end time.Time) ([]Datapoint, error) {
	f, err := os.Open(filepath.Join(Config.DBDir, date+".csv"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// read header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var datapoints []Datapoint

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return datapoints, err
		}

		timestamp, err := time.Parse(TIMESTAMP_FORMAT, record[0])
		if err != nil {
			return datapoints, err
		}
		temperature, _ := strconv.ParseFloat(record[2], 32)
		humidity, _ := strconv.ParseFloat(record[3], 32)
		pressure, _ := strconv.ParseFloat(record[4], 32)
		battery, _ := strconv.ParseUint(record[5], 10, 32)

		datapoint := Datapoint{
			Timestamp:   record[0],
			Station:     record[1],
			Temperature: float32(temperature),
			Humidity:    float32(humidity),
			Pressure:    float32(pressure),
			Battery:     uint(battery),
		}

		if datapoint.Station == station && timestamp.After(start) && timestamp.Before(end) {
			datapoints = append(datapoints, datapoint)
		}
	}

	return datapoints, nil
}

func FetchData(req *DataRequest) ([]Datapoint, error) {
	start, err := time.Parse(TIMESTAMP_FORMAT, req.TimeFrom)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(TIMESTAMP_FORMAT, req.TimeTo)
	if err != nil {
		return nil, err
	}

	if start.After(end) {
		return nil, fmt.Errorf("start timestamp must be before or equal to end timestamp")
	}

	data := []Datapoint{}
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		curData, err := ReadDataForStation(req.Station, current.Format(DATE_FORMAT), start, end)
		if err != nil {
			return nil, err
		}

		data = append(data, curData...)
	}

	return data, nil
}
