package metar

import (
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
		URL string  = "https://aviationweather.gov/cgi-bin/data/metar.php?ids=%s&hours=0&order=id%%2C-obs&sep=true&taf=%s"
		
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

	reports, err := parseResponse(resp, codes, tafOn)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

/* 
	Parses the response body, looking for METAR and, optionally, TAF phrases. 
	Errors returned relate to resp.Body reading problems.
*/
func parseResponse(resp *http.Response, codes []string, taf bool) ([]*Finding, error) {
	const (
		// Report signatures used for finding right content and ensuring the proper format
		METAR_SIG = "METAR"
		TAF_SIG   = "TAF"
	)

	stream, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stream), "\n\n")

	finds := make([]*Finding, len(codes))
	for i := range finds {
		finds[i] = &Finding{strings.ToUpper(codes[i]), "", "", true}
	}

	for i := 0; i < len(lines); i++ {
		f := pointerToFinding(finds, lines[i])
		
		if f == nil {
			continue
		}

		metar := strings.TrimRight(lines[i], "\n")

		if strings.HasPrefix(metar, METAR_SIG) {
			f.METAR = metar
		} else {
			f.METAR = METAR_SIG + " " + metar
		}

		if strings.HasSuffix(f.METAR, "$") {
			f.OK = false
		}

		if !taf {
			continue
		}

		if i < len(lines) - 1 {
			tafR := strings.TrimRight(lines[i + 1], "\n")

			if strings.HasPrefix(tafR, TAF_SIG) {
				f.TAF = tafR
			} else {
				f.TAF = TAF_SIG + " " + tafR
			}

			i++
		}
	}

	return finds, nil
}
