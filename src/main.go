package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zedran/weather-reports/src/metar"
)

func main() {
	log.SetFlags(0)

	// HTTP client timeout in seconds
	const CLIENT_TIMEOUT = 30

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

	client := http.Client{
		Timeout: CLIENT_TIMEOUT * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	fmt.Println(metar.GetReport(&client, cleanCodes, !*noTAF))
}
