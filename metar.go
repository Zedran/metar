package metar

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

/*
Sends request to the website and returns parsed results as a slice of Finding structures.
Errors returned are related to http package and parseResponse function.
*/
func GetReports(client *http.Client, codes []string, tafOn bool) ([]*Finding, error) {
	const (
		// Source URL
		URL string = "https://aviationweather.gov/api/data/metar?ids=%s&format=raw&taf=%s"

		// Character delimiting codes within the link
		CODES_DELIM = ","
	)

	var taf string

	if tafOn {
		taf = "true"
	} else {
		taf = "false"
	}

	resp, err := client.Get(fmt.Sprintf(URL, strings.Join(codes, CODES_DELIM), taf))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	stream, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseResponse(string(stream), tafOn)
}

/*
Parses the response body, looking for METAR and, optionally, TAF phrases.
Errors returned relate to resp.Body reading problems or unexpected response
format.
*/
func parseResponse(resp string, tafOn bool) ([]*Finding, error) {
	lines := strings.Split(resp, "\n")

	finds := make([]*Finding, 0)

	var (
		b       strings.Builder
		current *Finding
	)

	for _, ln := range lines {
		if strings.HasPrefix(ln, "METAR") {
			if current != nil {
				if tafOn {
					current.TAF = b.String()
				} else {

					current.METAR = b.String()
				}
				b.Reset()
			}
			current = new(Finding)
			current.Code = ln[6:10]
			finds = append(finds, current)
			b.WriteString(strings.TrimRight(ln, " ") + "\n")
			continue
		}

		if strings.HasPrefix(ln, "TAF") {
			if current == nil {
				return nil, errors.New("unexpected formatting")
			}
			current.METAR = b.String()
			b.Reset()
			b.WriteString(strings.TrimRight(ln, " ") + "\n")
			continue
		}

		b.WriteString(strings.TrimRight(ln, " ") + "\n")
	}

	if b.Len() > 0 {
		if tafOn {
			current.TAF = b.String()
		} else {
			current.METAR = b.String()
		}
	}

	for _, f := range finds {
		f.METAR = strings.TrimRight(f.METAR, " \n")
		f.TAF = strings.TrimRight(f.TAF, " \n")
		if !strings.HasSuffix(f.METAR, "$") {
			f.OK = true
		}
	}

	return finds, nil
}
