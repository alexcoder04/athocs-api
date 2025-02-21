package main

func GetStations() ([]Station, error) {
	return ReadCSV(Config.StationsIndex, StationFromCSV)
}

func AddStation(station string) error {
	stations, err := GetStations()
	if err != nil {
		return err
	}

	// check if already exists
	for _, st := range stations {
		if st.ID == station {
			return nil
		}
	}

	return WriteToCSV(Config.StationsIndex, Station{
		ID:     station,
		Name:   station,
		Active: 1,
	}.ToCSVRow())
}
