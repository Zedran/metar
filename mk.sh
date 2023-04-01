#!/usr/bin/env bash

source_path=./src
output_path=./build/metar

go build -o $output_path -trimpath -ldflags "-s -w" $source_path
