package metar

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseResponse(t *testing.T) {
	type testCase struct {
		InputFile string     `json:"input_file"`
		TafOn     bool       `json:"taf_on"`
		Expected  []*Finding `json:"expected"`
	}

	var cases []testCase

	testData, err := os.ReadFile("testdata/cases.json")
	if err != nil {
		t.Fatal("failed to load cases.json")
	}

	if err := json.Unmarshal(testData, &cases); err != nil {
		t.Fatal("failed to unmarshal cases.json")
	}

	for _, c := range cases {
		input, err := os.ReadFile(filepath.Join("testdata", c.InputFile))
		if err != nil {
			t.Fatalf("failed to load %s", c.InputFile)
		}

		out, err := parseResponse(string(input), c.TafOn)
		if err != nil {
			t.Errorf("parseResponse returned an unexpected error: %v", err)
		}

		if len(c.Expected) != len(out) {
			t.Fatalf("length mismatch, out == '%v'", out)
		}

		for _, ef := range c.Expected {
			found := false

			for _, of := range out {
				if ef.Code == of.Code && ef.METAR == of.METAR && ef.TAF == of.TAF && ef.OK == of.OK {
					found = true
					break
				}
			}
			if !found {
				t.Logf("%s: %+v was not found among:", c.InputFile, ef)
				for _, f := range out {
					t.Logf("%+v", *f)
				}
				t.FailNow()
			}
		}

	}
}
