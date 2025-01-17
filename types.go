package main

type Datapoint struct {
    Timestamp string `form:"timestamp" json:"timestamp"`
    Temperature float32 `form:"temperature" json:"temperature"`
    Humidity float32 `form:"humidity" json:"humidity"`
    Pressure float32 `form:"pressure" json:"pressure"`
    Station string `form:"station" json:"station"`
    Battery uint `form:"battery" json:"battery"`
}

