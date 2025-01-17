package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

var DatabaseFolder = GetDatabaseFolder()

func GetDatabaseFolder() string {
    dir := os.Getenv("ATHOCS_DATABASE_DIR")
    if dir == "" {
        dir = "./data"
    }

    if _, err := os.Stat(dir); os.IsNotExist(err) {
        err := os.MkdirAll(dir, 0700)
        if err != nil {
            panic("data directory does not exist and creating it failed")
        }
    } else if err != nil {
        panic("unknown error while checking for data directory")
    }

    return dir
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
    parsedDate, err := time.Parse("2006-01-02_15:04:05", timestamp)
    if err != nil {
        return "", err
    }

    dateString := parsedDate.Format("2006-01-02")
    return filepath.Join(DatabaseFolder, dateString + ".csv"), nil
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
        fmt.Println(err.Error())
    }
    return err
}

