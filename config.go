package main

import (
    "os"

    "github.com/alexcoder04/ef"
)

type AthocsConfig struct {
    BaseDir string
    DBDir string
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
        BaseDir: baseDir.Path,
        DBDir: dbDir.Path,
        Port: port,
    }
}

