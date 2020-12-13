package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const URL string = "https://www.aviationweather.gov/metar/data?ids=%s&format=raw&hours=0&taf=%s&layout=on"

var (
	errReportNotFound = errors.New("Report not found. The airfield code may be invalid.")

	client = http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
)

func fatal(message string) {
	fmt.Println(message)
	flag.PrintDefaults()
	os.Exit(1)
}

func getReport(code string, tafOn bool) string {
	var taf string
	
	if tafOn {
		taf = "on"
	} else {
		taf = "off"
	}

	resp, err := client.Get(fmt.Sprintf(URL, code, taf))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	report, err := parseReport(resp, code, tafOn)
	if err != nil {
		log.Fatal(err)
	}

	return report
}

func parseReport(resp *http.Response, code string, taf bool) (string, error) {
	const METAR_NF_PHRASE string = "No METAR found for"
	
	page := html.NewTokenizer(resp.Body)

	output := make([]string, 0)
	for page.Next() != html.ErrorToken {
		token := page.Token()

		if token.Type == html.TextToken {
			str := token.String()
			if strings.Contains(str, code) {
				output = append(output, str)
			}
		}
	}

	var err error
	if len(output) == 0 || strings.Contains(output[0], METAR_NF_PHRASE) {
		err = errReportNotFound
	}

	return "METAR " + strings.Join(output, "\n"), err
}

func main() {
	log.SetFlags(0)

	action := flag.String("a", "m", "action:\n    m - get METAR\n    l - display download link\n")
	code   := flag.String("c", "", "a 4-letter ICAO airport code")
	noTAF  := flag.Bool("notaf", false, "do not get TAF report")

	flag.Parse()

	switch strings.ToUpper(*action) {
	case "L":
		fmt.Printf(URL, "<CODE>", "<on/off>")
	case "M":
		if len(*code) != 4 {
			fatal("ICAO code not specified or of incorrect length.")
		}
		fmt.Println(getReport(strings.ToUpper(*code), !*noTAF))
	default:
		fatal("Unknown action flag.")
	}
}
