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

const (
	// Source URL
	URL string    = "https://www.aviationweather.gov/metar/data?ids=%s&format=raw&hours=0&taf=%s&layout=off"
	
	// HTTP request timeout time
	TIMEOUT_TIME  = 30
	
	// Airport ICAO code length
	ICAO_CODE_LEN =  4
)

var (
	errReportNotFound = errors.New("Report not found. The airfield code may be invalid.")

	client = http.Client{
		Timeout: TIMEOUT_TIME * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
)

// Displays feedback message, prints flags list and exits the program with 1 status.
func fatal(message string) {
	fmt.Println(message)
	flag.PrintDefaults()
	os.Exit(1)
}

// Gets report from website and returns parsed result.
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

// Parses the response body, looking for METAR and, optionally, TAF phrases. Returns error
// if the report for the given ICAO code was not found.
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

	action := flag.String("a", "m", "action:\n    m (metar) - get METAR\n    l (link)  - display download link\n")
	code   := flag.String("c", "", "a 4-letter ICAO airport code, you can specify it without a flag")
	noTAF  := flag.Bool("notaf", false, "do not get TAF report")

	flag.Parse()

	switch strings.ToUpper(*action) {
	case "L", "LINK":
		fmt.Printf(URL, "<CODE>", "<on/off>")
	case "M", "METAR":
		if len(*code) != ICAO_CODE_LEN {
			if *code = flag.Arg(0); len(*code) != ICAO_CODE_LEN {
				fatal("ICAO code not specified or of incorrect length.")
			}
		}
		fmt.Println(getReport(strings.ToUpper(*code), !*noTAF))
	default:
		fatal(fmt.Sprintf("Unknown action flag value '%s'.", *action))
	}
}
