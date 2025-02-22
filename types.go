package main

import (
	"fmt"
	"strconv"
	"time"
)

type DataRequest struct {
	Station  string `query:"station"`
	TimeFrom string `query:"time_from"`
	TimeTo   string `query:"time_to"`
}

type Station struct {
	ID     string `form:"id" json:"id"`
	Name   string `form:"name" json:"name"`
	Active uint64 `form:"active" json:"active"`
}

func (st Station) ToCSVRow() []string {
	return []string{
		st.ID,
		st.Name,
		fmt.Sprintf("%d", st.Active),
	}
}

func (st Station) Header() []string {
	return []string{"id", "name", "active"}
}

func (st Station) Hidden() bool {
	return st.Active == 0
}

func StationFromCSV(row []string) (Station, error) {
	active, err := strconv.ParseUint(row[2], 10, 32)
	return Station{
		ID:     row[0],
		Name:   row[1],
		Active: active,
	}, err
}

type Datapoint struct {
	Timestamp   string  `form:"timestamp" json:"timestamp"`
	Station     string  `form:"station" json:"station"`
	Temperature float64 `form:"temperature" json:"temperature"`
	Humidity    float64 `form:"humidity" json:"humidity"`
	Pressure    float64 `form:"pressure" json:"pressure"`
	Battery     uint64  `form:"battery" json:"battery"`
}

func (d Datapoint) ToCSVRow() []string {
	return []string{
		d.Timestamp,
		d.Station,
		fmt.Sprintf("%.2f", d.Temperature),
		fmt.Sprintf("%.2f", d.Humidity),
		fmt.Sprintf("%.2f", d.Pressure),
		fmt.Sprintf("%d", d.Battery),
	}
}

func (d Datapoint) Header() []string {
	return []string{"timestamp", "station", "temperature", "humidity", "pressure", "battery"}
}

func (d Datapoint) Hidden() bool {
	return false
}

func DatapointFromCSV(row []string) (Datapoint, error) {
	// validate that timestamp is valid
	_, err := time.Parse(Config.TimestampFormat, row[0])
	if err != nil {
		return Datapoint{}, err
	}

	temperature, _ := strconv.ParseFloat(row[2], 32)
	humidity, _ := strconv.ParseFloat(row[3], 32)
	pressure, _ := strconv.ParseFloat(row[4], 32)
	battery, _ := strconv.ParseUint(row[5], 10, 32)

	return Datapoint{
		Timestamp:   row[0],
		Station:     row[1],
		Temperature: temperature,
		Humidity:    humidity,
		Pressure:    pressure,
		Battery:     battery,
	}, nil
}

type CSVRowData interface {
	Station | Datapoint
	Header() []string
	ToCSVRow() []string
	Hidden() bool
}
