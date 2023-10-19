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

/* Removes duplicates and changes all letters in codes to upper case. */
func PrepareCodesList(codes []string) []string {
	processed := removeDuplicates(codes)

	for i := range processed {
		processed[i] = strings.ToUpper(processed[i])
	}
	
	return processed
}

/* Returns the slice without duplicate strings. */
func removeDuplicates(codes []string) []string {
	clean := make([]string, 0, len(codes))

	for i := range codes {
		if !contains(clean, codes[i]) {
			clean = append(clean, codes[i])
		}
	}

	return clean
}
