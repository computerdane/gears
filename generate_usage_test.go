package gears

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func makeWhiteSpaceVisible(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, " ", "_"), "\n", "\\n\n")
}

func TestFprintUsageWithWidth(t *testing.T) {
	tests_reset()

	target := `--long / -l
    A flag with a super duper long description. Like, this is a very long 
    description and is totally overwhelming the user. We really need to stop 
    making things so long and complicated guys. The poor users can't handle it! 
    (default: long)

--name / -n
    The person we want to greet (default: john)

--zzz
    An argument with no shorthand! (default: false)
`

	Add(&Flag{
		Name:         "name",
		ValueType:    "string",
		DefaultValue: "john",
		Description:  "The person we want to greet",
		Shorthand:    "n",
	})
	Add(&Flag{
		Name:         "long",
		ValueType:    "string",
		DefaultValue: "long",
		Description:  "A flag with a super duper long description. Like, this is a very long description and is totally overwhelming the user. We really need to stop making things so long and complicated guys. The poor users can't handle it!",
		Shorthand:    "l",
	})
	Add(&Flag{
		Name:        "zzz",
		ValueType:   "bool",
		Description: "An argument with no shorthand!",
	})
	Add(&Flag{
		Name:             "excluded",
		ValueType:        "bool",
		Description:      "This flag is excluded from the usage string",
		ExcludeFromUsage: true,
	})

	var w bytes.Buffer
	FprintUsageWithWidth(&w, 80)

	output := w.String()
	if output != target {
		fmt.Println(makeWhiteSpaceVisible(output))
		fmt.Println(makeWhiteSpaceVisible(target))
		t.Error("usage string was not as expected")
	}
}
