package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Zedran/weather-reports/src/metar"
)

func main() {
	log.SetFlags(0)

	noTAF := flag.Bool("notaf", false, "do not get TAF report")

	flag.Parse()

	codes := flag.Args()

	if len(codes) == 0 {
		log.Fatal("No code specified.")
	}

	cleanCodes := metar.PrepareCodes(codes)

	if len(cleanCodes) == 0 {
		log.Fatal("No valid code specified.")
	}

	fmt.Println(metar.GetReport(cleanCodes, !*noTAF))
}
