package metar

/* Finding structure holds the extracted airport data. */
type Finding struct {
	// Airport Code
	Code string

	METAR string
	TAF   string

	// Sensor's health: false indicates that sensor at the airport
	// requires maintenance and some data may be inaccurate
	OK bool
}
