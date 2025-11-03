package metar

import "strings"

/* Returns true if code is inside slice s. */
func contains(s []string, code string) bool {
	for i := range s {
		if s[i] == code {
			return true
		}
	}
	return false
}

/* Removes duplicates, invalid codes and changes every code to upper case. */
func PrepareCodes(codes ...string) []string {
	const ICAO_CODE_LEN = 4

	clean := make([]string, 0, len(codes))

	for i := range codes {
		codes[i] = strings.ToUpper(codes[i])
	}

	for i := range codes {
		if !contains(clean, codes[i]) && len(codes[i]) == ICAO_CODE_LEN {
			clean = append(clean, codes[i])
		}
	}

	return clean
}
