package main

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

/* Returns a pointer to the Finding that holds the airport code (Finding.Code) present within line. */
func PointerToFinding(finds []*Finding, line string) *Finding {
	for i := range finds {
		if strings.Contains(line, finds[i].Code) {
			return finds[i]
		}
	}
	return nil
}

/* Removes duplicates, invalid codes and changes every code to upper case. */
func PrepareCodes(codes []string) []string {
	const ICAO_CODE_LEN = 4

	clean := make([]string, 0, len(codes))

	for i := range codes {
		if !contains(clean, codes[i]) && len(codes[i]) == ICAO_CODE_LEN {
			clean = append(clean, codes[i])
		}
	}

	for i := range clean {
		clean[i] = strings.ToUpper(clean[i])
	}

	return clean
}
