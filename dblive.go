package main

import (
	"path/filepath"
	"time"
)

type LiveDatabase struct {
	Data map[string]Datapoint
}

var LiveDB = LiveDatabase{
	Data: map[string]Datapoint{},
}

func (livedb *LiveDatabase) Init() error {
	// we only care about last 24 hours
	now := time.Now()
	limit := now.Add(-24 * time.Hour)

	// run through dates in reverse order
	for current := now; !current.Before(limit); current = current.Add(-24 * time.Hour) {
		data, err := ReadCSV(filepath.Join(Config.DBDir, current.Format(Config.DateFormat)+".csv"), DatapointFromCSV)
		if err != nil {
			return err
		}

		foundStations := map[string]bool{}
		for i := len(data) - 1; i > 0; i-- {
			// we already have the latest values
			if _, exists := foundStations[data[i].Station]; exists {
				continue
			}

			// add value
			livedb.Data[data[i].Station] = data[i]
			foundStations[data[i].Station] = true
		}
	}

	return nil
}

func (livedb *LiveDatabase) Set(dp Datapoint) {
	livedb.Data[dp.Station] = dp
}

func (livedb *LiveDatabase) Get(station string) Datapoint {
	return livedb.Data[station]
}
