#!/usr/bin/env bash

go build -o ./build/metar -ldflags "-s -w" ./src
