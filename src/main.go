package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Zedran/weather-reports/src/metar"
)

const (
	// This phrase is returned when the code is not assigned to any airport
	NF_PHRASE  = "No %s found for %s"

	// Output delimiter used when more than one airport code is specified
	OUT_DELIM  = "\n\n---------------------------------------\n\n"

	// Output format, mirrors the website's way of displaying reports
	OUT_FORMAT = "%s\n\n%s"
)

/* Returns the string containing the formatted report. */
func FindingToString(f *metar.Finding, taf bool) string {
	if len(f.METAR) == 0 {
		f.METAR = fmt.Sprintf(NF_PHRASE, "METAR", f.Code)
	}

	if !taf {
		return f.METAR
	}

	if len(f.TAF) == 0 {
		f.TAF = fmt.Sprintf(NF_PHRASE, "TAF", f.Code)
	}

	return fmt.Sprintf(OUT_FORMAT, f.METAR, f.TAF)
}

/* Calls FindingToString for every Finding. */
func PrintFindings(f []*metar.Finding, taf bool) {
	var b strings.Builder

	fmt.Fprint(&b, "\n")

	for i := range f {
		fmt.Fprint(&b, FindingToString(f[i], taf))

		if i < len(f) - 1 {
			fmt.Fprint(&b, OUT_DELIM)
		}
	}

	fmt.Fprint(&b, "\n")

	fmt.Println(b.String())
}

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

	findings, err := metar.GetReports(&client, cleanCodes, !*noTAF)
	if err != nil {
		log.Fatal(err)
	}

	PrintFindings(findings, !*noTAF)
}
