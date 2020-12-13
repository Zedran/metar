package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	URL string       = "https://www.aviationweather.gov/metar/data?ids=%s&format=raw&hours=0&taf=%s&layout=on"

	HELP_MSG string  = "Proper format: 'metar <code> [notaf]' or 'metar <link>'"

	MIN_ARGC         = 2  // metar code
	MAX_ARGC         = 3  // metar code notaf

	NOTAF_ARG string = "NOTAF"
	LINK_ARG  string = "LINK"
)

var (
	errInvalidArguments = errors.New("Invalid arguments. An ICAO airport code must be specified. Add 'taf' for TAF report.")
	errReportNotFound   = errors.New("Report not found. The airfield code may be invalid.")

	client = http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
)

func getReport(code string, tafOn bool) (string, error) {
	var taf string
	
	if tafOn {
		taf = "on"
	} else {
		taf = "off"
	}

	resp, err := client.Get(fmt.Sprintf(URL, code, taf))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return parseReport(resp, code, tafOn)
}

func parseReport(resp *http.Response, code string, taf bool) (string, error) {
	page := html.NewTokenizer(resp.Body)

	var (
		str string
		output     = make([]string, 0)
		p bool = false
	)

	for page.Next() != html.ErrorToken {
		token := page.Token()

		if token.Type == html.TextToken {
			str = token.String()
			if p {
				print(str)
			}
			if strings.Contains(str, code) {
				output = append(output, str)
			}
		} else if token.Type == html.CommentToken {
			if strings.Contains(token.String(), "Data starts here") {
				p = true
			} else if strings.Contains(token.String(), "Data ends here") {
				p = false
			}
		}
	}

	var err error
	if len(output) == 0 {
		err = errReportNotFound
	}

	return "METAR " + strings.Join(output, "\n"), err
}

func main() {
	log.SetFlags(0)

	taf := true

	switch len(os.Args) {
	case 2:
		if strings.ToUpper(os.Args[1]) == LINK_ARG {
			fmt.Printf(URL, "<CODE>", "<on/off>")
			os.Exit(0)
		}
	case 3:
		if strings.ToUpper(os.Args[2]) == NOTAF_ARG {
			taf = false
			break
		}
		fallthrough
	default:
		log.Fatal(errInvalidArguments)
	}

	report, err := getReport(strings.ToUpper(os.Args[1]), taf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(report)
}
