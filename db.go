package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	TIMESTAMP_FORMAT            = "2006-01-02_15:04:05"
	DATE_FORMAT                 = "2006-01-02"
	DEFAULT_DATA_INTERVAL_HOURS = 24 * 7
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

func GetStationsList() ([]Station, error) {
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

	var stations []Station

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return stations, err
		}

		// TODO auto deactivate stations
		active, _ := strconv.ParseUint(record[2], 10, 32)
		st := Station{
			ID:     record[0],
			Name:   record[1],
			Active: uint(active),
		}
		stations = append(stations, st)
	}

	return stations, nil
}

func AddStation(station string) error {
	stations, err := GetStationsList()
	if err != nil {
		return err
	}

	// check if already exists
	for _, st := range stations {
		if st.ID == station {
			return nil
		}
	}

	f, err := os.OpenFile(filepath.Join(Config.DBDir, "stations-index.csv"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf(
		"%s,%s,1\n",
		station,
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
		if os.IsNotExist(err) {
			return nil, nil
		}
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
	var end time.Time
	var err error
	if req.TimeTo == "" {
		end = time.Now()
	} else {
		end, err = time.Parse(TIMESTAMP_FORMAT, req.TimeTo)
		if err != nil {
			return nil, err
		}
	}

	var start time.Time
	if req.TimeFrom == "" {
		start = end.Add(-DEFAULT_DATA_INTERVAL_HOURS * time.Hour)
	} else {
		start, err = time.Parse(TIMESTAMP_FORMAT, req.TimeFrom)
		if err != nil {
			return nil, err
		}
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
