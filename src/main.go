package main

import (
	"flag"
	"fmt"
	"log"
)

// Airport ICAO code length
const ICAO_CODE_LEN = 4

/* Returns true if code is inside slice s. */
func Contains(s []string, code string) bool {
	for i := range s {
		if s[i] == code {
			return true
		}
	}
	return false
}

/* Rewrites the passed slice, omitting duplicate values. */
func RemoveDuplicates(codes []string) []string {
	clean := make([]string, 0, len(codes))

	for i := range codes {
		if !Contains(clean, codes[i]) {
			clean = append(clean, codes[i])
		}
	}

	return clean
}

func main() {
	log.SetFlags(0)

	noTAF := flag.Bool("notaf", false, "do not get TAF report")

	flag.Parse()

	codes := flag.Args()

	if len(codes) == 0 {
		log.Fatal("No code specified.")
	}

	for i := range codes {
		if len(codes[i]) != ICAO_CODE_LEN {
			log.Fatal(fmt.Sprintf("Improper code format: '%s'", codes[i]))
		}
	}

	fmt.Println(GetReport(RemoveDuplicates(codes), !*noTAF))
}
