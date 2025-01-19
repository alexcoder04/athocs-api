#!/bin/sh

GOOS=linux GOARCH=arm GOARM=6 go build -o athocs-api \
    && rsync -vu athocs-api "192.168.0.86:/home/alex/athocs/api/" \
    && rm athocs-api

