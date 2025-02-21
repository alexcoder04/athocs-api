package main

import (
	"os"
	"path/filepath"

	"github.com/alexcoder04/ef"
)

type AthocsConfig struct {
	TimestampFormat string

	BaseDir       string
	DBDir         string
	StationsIndex string

	Port string
}

var Config = GetConfig()

func GetConfig() AthocsConfig {
	baseDir := ef.NewFile(os.Getenv("ATHOCS_BASE_DIR"))
	isDir, err := baseDir.IsDir()
	if err != nil {
		panic("unknown error while checking for athocs dir")
	}
	if !isDir {
		panic("athocs dir non existent")
	}

	dbDir := ef.NewFile(baseDir.Path, "data")
	isDir, err = dbDir.IsDir()
	if err != nil {
		panic("unknown error while checking for database dir")
	}
	if !isDir {
		panic("database dir non existent")
	}

	port := os.Getenv("ATHOCS_PORT")
	if port == "" {
		port = "1111"
	}

	return AthocsConfig{
		TimestampFormat: "2006-01-02_15:04:05",
		BaseDir:         baseDir.Path,
		DBDir:           dbDir.Path,
		StationsIndex:   filepath.Join(dbDir.Path, "stations-index.csv"),
		Port:            port,
	}
}
