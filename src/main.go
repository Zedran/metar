package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	log.SetFlags(0)

	noTAF := flag.Bool("notaf", false, "do not get TAF report")

	flag.Parse()

	codes := flag.Args()

	if len(codes) == 0 {
		log.Fatal("No code specified.")
	}

	cleanCodes := PrepareCodes(codes)

	if len(cleanCodes) == 0 {
		log.Fatal("No valid code specified.")
	}

	fmt.Println(GetReport(cleanCodes, !*noTAF))
}
