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

	// Character delimiting codes within the link
	CODES_DELIM   = "+"
	
	// HTTP request timeout time
	TIMEOUT_TIME  = 30
	
	// Airport ICAO code length
	ICAO_CODE_LEN =  4

	// Report signatures used for finding right content and ensuring the proper format
	METAR_SIG     = "METAR"
	TAF_SIG       = "TAF"

	// '<br/>' and '&nbsp;' substitutes used for output formatting
	BR            = "\n"
	NBSP          = " "

	// Output delimiter used when more than one airport code is specified
	OUT_DELIM     = "\n\n---------------------------------------\n\n"
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

/* Extracts reports from the website and returns parsed result. */
func GetReport(codes []string, tafOn bool) string {
	var taf string
	
	if tafOn {
		taf = "on"
	} else {
		taf = "off"
	}

	resp, err := client.Get(fmt.Sprintf(URL, strings.Join(codes, CODES_DELIM), taf))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	report, err := parseReport(resp, codes, tafOn)
	if err != nil {
		log.Fatal(err)
	}

	return report
}

/* Parses the response body, looking for METAR and, optionally, TAF phrases.
 * Errors returned relate to resp.Body reading problems.
 */
func parseReport(resp *http.Response, codes []string, taf bool) (string, error) {
	stream, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := exp.FindAllSubmatch(stream, -1)

	finds := make([]*Finding, len(codes))
	for i := range finds {
		finds[i] = NewFinding(strings.ToUpper(codes[i]))
	}

	for i := range matches {
		str := string(matches[i][1])
		
		for i := range finds {
			if strings.Contains(str, finds[i].Code) {
				if strings.Contains(str, TAF_SIG) {
					finds[i].TAF = strings.ReplaceAll(strings.ReplaceAll(str, "<br/>", BR), "&nbsp;", NBSP)
				} else {
					finds[i].METAR = METAR_SIG + NBSP + str
				}
				break
			}
		}
	}

	reports := make([]string, len(finds))

	for i := range finds {
		reports[i] = finds[i].ToString(taf)
	}

	return strings.Join(reports, OUT_DELIM), nil
}

/* Returns true if code is inside slice s. */
func Contains(s []string, code string) bool {
	for i := range s {
		if s[i] == code {
			return true
		}
	}
	return false
}

/* Rewrites the passed slice, omitting duplicate values. */
func RemoveDuplicates(codes []string) []string {
	clean := make([]string, 0, len(codes))

	for i := range codes {
		if !Contains(clean, codes[i]) {
			clean = append(clean, codes[i])
		}
	}

	return clean
}

func main() {
	log.SetFlags(0)

	noTAF := flag.Bool("notaf", false, "do not get TAF report")

	flag.Parse()

	codes := flag.Args()

	if len(codes) == 0 {
		fatal("No code specified.")
	}

	for i := range codes {
		if len(codes[i]) != ICAO_CODE_LEN {
			fatal(fmt.Sprintf("Improper code format: %s", codes[i]))
		}
	}

	fmt.Println(GetReport(RemoveDuplicates(codes), !*noTAF))
}
