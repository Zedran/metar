@ECHO OFF

set source_path=./src
set output_path=./build/metar.exe

go build -o %output_path% -trimpath -ldflags "-s -w" %source_path%
