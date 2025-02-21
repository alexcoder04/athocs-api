package main

import "fmt"

type DataRequest struct {
	Station  string `query:"station"`
	TimeFrom string `query:"time_from"`
	TimeTo   string `query:"time_to"`
}

type Station struct {
	ID     string `form:"id" json:"id"`
	Name   string `form:"name" json:"name"`
	Active uint   `form:"active" json:"active"`
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

type Datapoint struct {
	Timestamp   string  `form:"timestamp" json:"timestamp"`
	Station     string  `form:"station" json:"station"`
	Temperature float32 `form:"temperature" json:"temperature"`
	Humidity    float32 `form:"humidity" json:"humidity"`
	Pressure    float32 `form:"pressure" json:"pressure"`
	Battery     uint    `form:"battery" json:"battery"`
}

func (d Datapoint) ToCSVRow() []string {
	return []string{
		d.Timestamp,
		d.Station,
		fmt.Sprintf("%f", d.Temperature),
		fmt.Sprintf("%f", d.Humidity),
		fmt.Sprintf("%f", d.Pressure),
		fmt.Sprintf("%d", d.Battery),
	}
}

func (d Datapoint) Header() []string {
	return []string{"timestamp", "station", "temperature", "humidity", "pressure", "battery"}
}

func (d Datapoint) Hidden() bool {
	return false
}

type CSVRowData interface {
	Station | Datapoint
	Header() []string
	ToCSVRow() []string
	Hidden() bool
}
