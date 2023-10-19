package main

import (
	"flag"
	"fmt"
	"log"
)

// Airport ICAO code length
const ICAO_CODE_LEN = 4

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

	fmt.Println(GetReport(PrepareCodesList(codes), !*noTAF))
}
