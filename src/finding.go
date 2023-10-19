package main

import "fmt"

const (
	// Output format, mirrors the website's way of displaying reports
	OUT_FORMAT = "%s\n\n%s"

	// This phrase is returned when the code is not assigned to any airport
	NF_PHRASE  = "No %s found for %s"
)

/* Finding structure holds the extracted airport data. */
type Finding struct {
	// Airport Code
	Code  string
	
	METAR string
	TAF   string
}

/* Returns the string containing the formatted report. */
func (f *Finding) ToString(taf bool) string {
	if len(f.METAR) == 0 {
		f.METAR = fmt.Sprintf(NF_PHRASE, METAR_SIG, f.Code)
	}

	if !taf {
		return f.METAR
	}

	if len(f.TAF) == 0 {
		f.TAF = fmt.Sprintf(NF_PHRASE, TAF_SIG, f.Code)
	}

	return fmt.Sprintf(OUT_FORMAT, f.METAR, f.TAF)
}

/* Creates new Finding struct and assigns the airport code to it. */
func NewFinding(code string) *Finding {
	return &Finding{code, "", ""}
}
