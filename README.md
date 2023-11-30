# METAR

A simple package providing the ability to download and parse the most recent aerodrome weather reports (METAR and TAF). It utilizes [Aviation Weather Center](https://aviationweather.gov) as data source.

## Usage

### Prepare ICAO codes before query

Remove duplicates and improperly formatted codes.

```Go
func PrepareCodes(codes ...string) []string
```

### Get weather reports from Aviation Weather Center

```Go
func GetReports(client *http.Client, codes []string, tafOn bool) ([]*Finding, error)
```

## Resources

* [ICAO airport code](https://en.wikipedia.org/wiki/ICAO_airport_code)
* [Aviation Weather Center](https://aviationweather.gov)
* [Aviation Weather Center: Disclaimer](https://www.weather.gov/disclaimer#public-note-of-appropriate-use)

## License

This software is available under MIT License.
