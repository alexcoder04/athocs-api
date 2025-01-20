#!/bin/sh

. ~/.profile \
    && cd "$ATHOCS_BASE_DIR" \
    && exec api/athocs-api

