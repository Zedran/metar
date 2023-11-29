package metar

/* Finding structure holds the extracted airport data. */
type Finding struct {
	// Airport Code
	Code  string
	
	METAR string
	TAF   string
}

/* Creates new Finding struct and assigns the airport code to it. */
func NewFinding(code string) *Finding {
	return &Finding{code, "", ""}
}
