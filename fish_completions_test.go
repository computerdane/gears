package gears

import (
	"log"
	"strings"
	"testing"
)

func TestFishCompletions(t *testing.T) {
	tests_reset()

	if err := Add(&Flag{Name: "my-str", Shorthand: "s", ValueType: "string", DefaultValue: "", Description: "My string description"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "my-bool", Shorthand: "b", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "my-floats", ValueType: "floats", DefaultValue: []float64{}}); err != nil {
		log.Fatal(err)
	}

	expectedLines := []string{
		`complete -c "my-cmd" -l "my-str" -s "s" -d "My string description"`,
		`complete -c "my-cmd" -l "my-bool" -s "b"`,
		`complete -c "my-cmd" -l "my-floats"`,
	}

	completions := FishCompletions("my-cmd")
	for _, line := range strings.Split(strings.TrimSpace(completions), "\n") {
		found := false
		for _, expected := range expectedLines {
			if line == expected {
				found = true
			}
		}
		if !found {
			t.Errorf("completion not found in expected list: %s", line)
		}
	}
}
