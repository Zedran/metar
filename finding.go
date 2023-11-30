package metar

/* Finding structure holds the extracted airport data. */
type Finding struct {
	// Airport Code
	Code  string
	
	METAR string
	TAF   string
}
