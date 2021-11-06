package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	// Source URL
	URL string    = "https://www.aviationweather.gov/metar/data?ids=%s&format=raw&hours=0&taf=%s&layout=off"

	// Character delimiting codes within the link
	CODES_DELIM   = "+"

	// Report signatures used for finding right content and ensuring the proper format
	METAR_SIG     = "METAR"
	TAF_SIG       = "TAF"

	// '<br/>' and '&nbsp;' substitutes used for output formatting
	BR            = "\n"
	NBSP          = " "

	// Output delimiter used when more than one airport code is specified
	OUT_DELIM     = "\n\n---------------------------------------\n\n"

	// HTTP request timeout time
	TIMEOUT_TIME  = 30
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

	report, err := parseResponse(resp, codes, tafOn)
	if err != nil {
		log.Fatal(err)
	}

	return report
}

/* Parses the response body, looking for METAR and, optionally, TAF phrases.
 * Errors returned relate to resp.Body reading problems.
 */
func parseResponse(resp *http.Response, codes []string, taf bool) (string, error) {
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
				if taf && len(finds[i].METAR) > 0 {
					// TAF always comes second, so it can be assumed that if len(f.METAR) > 0, the report is TAF
					finds[i].TAF = strings.ReplaceAll(strings.ReplaceAll(str, "<br/>", BR), "&nbsp;", NBSP)

					if !strings.Contains(finds[i].TAF, TAF_SIG) {
						// Since american reports do not have TAF signature in website's code, it is appended
						finds[i].TAF = TAF_SIG + NBSP + finds[i].TAF
					}
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
