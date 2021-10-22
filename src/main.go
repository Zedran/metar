package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	// Source URL
	URL string    = "https://www.aviationweather.gov/metar/data?ids=%s&format=raw&hours=0&taf=%s&layout=off"
	
	// HTTP request timeout time
	TIMEOUT_TIME  = 30
	
	// Airport ICAO code length
	ICAO_CODE_LEN =  4

	// Report signatures used for finding right content and ensuring the proper format
	METAR_SIG     = "METAR"
	TAF_SIG       = "TAF"
	
	// This phrase is returned when code is not assigned to any airport
	NF_PHRASE     = "No report found for %s"

	// '<br/>' and '&nbsp;' substitutes used for output formatting
	BR            = "\n"
	NBSP          = " "

	// output format, mirrors the website's way of displaying reports
	OUT_FORMAT    = "%s\n\n%s"
)

var (
	client = http.Client{
		Timeout: TIMEOUT_TIME * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	// Regexp for finding the content of the '<code>' tags
	exp = regexp.MustCompile("<code>(.*)</code>")
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
	stream, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := exp.FindAllSubmatch(stream, -1)

	if len(matches) == 0 {
		return fmt.Sprintf(NF_PHRASE, code), nil
	}

	var metarContent, tafContent string

	for i := range matches {
		str := string(matches[i][1])

		if strings.Contains(str, TAF_SIG) {
			tafContent = strings.ReplaceAll(strings.ReplaceAll(str, "<br/>", BR), "&nbsp;", NBSP)
		} else {
			metarContent = METAR_SIG + NBSP + string(matches[0][1])
		}
	}

	return fmt.Sprintf(OUT_FORMAT, metarContent, tafContent), nil
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
