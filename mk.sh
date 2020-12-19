#!/usr/bin/env bash

go build -o ./build/metar -trimpath -ldflags "-s -w" ./src
