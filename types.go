package main

type Datapoint struct {
    Timestamp string `form:"timestamp" json:"timestamp"`
    Station string `form:"station" json:"station"`
    Temperature float32 `form:"temperature" json:"temperature"`
    Humidity float32 `form:"humidity" json:"humidity"`
    Pressure float32 `form:"pressure" json:"pressure"`
    Battery uint `form:"battery" json:"battery"`
}

type DataRequest struct {
    Station string `query:"station"`
    TimeFrom string `query:"time_from"`
    TimeTo string `query:"time_to"`
}

