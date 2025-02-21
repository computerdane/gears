package gears

import (
	"fmt"
	"io"
	"os"
	"slices"
)

func appendDefaultValue(description *string, value any) {
	*description += fmt.Sprintf(" (default: %v)", value)
}

func FprintUsageWithWidth(w io.Writer, width int) {
	indent := "    "
	maxDescWidth := width - len(indent)

	names := make([]string, len(flags))
	i := 0
	for _, flag := range flags {
		names[i] = flag.Name
		i++
	}
	slices.Sort(names)

	for n, name := range names {
		flag := flags[name]
		if flag.ExcludeFromUsage || flag.Description == "" {
			continue
		}

		fmt.Fprintf(w, "--%s", flag.Name)
		if flag.Shorthand != "" {
			fmt.Fprintf(w, " / -%s", flag.Shorthand)
		}
		fmt.Fprintln(w)

		desc := flag.Description
		switch flag.ValueType {
		case "bool":
			appendDefaultValue(&desc, "false")
		case "float":
			appendDefaultValue(&desc, flag.DefaultValue)
		case "int":
			appendDefaultValue(&desc, flag.DefaultValue)
		case "string":
			if flag.DefaultValue != "" {
				appendDefaultValue(&desc, flag.DefaultValue)
			}
		case "floats":
			value := flag.DefaultValue.([]float64)
			if len(value) != 0 {
				appendDefaultValue(&desc, value)
			}
		case "ints":
			value := flag.DefaultValue.([]int)
			if len(value) != 0 {
				appendDefaultValue(&desc, value)
			}
		case "strings":
			value := flag.DefaultValue.([]string)
			if len(value) != 0 {
				appendDefaultValue(&desc, value)
			}
		}

		for l := 0; l < len(desc); {
			remaining := len(desc) - l
			maxWrappedWidth := min(maxDescWidth, remaining)
			wrappedWidth := maxWrappedWidth
			for remaining > maxDescWidth && desc[l+wrappedWidth-1] != ' ' {
				if wrappedWidth == 0 {
					// Could not find a space, just break up a word
					wrappedWidth = maxWrappedWidth
					break
				}
				wrappedWidth--
			}
			fmt.Fprintf(w, "%s%s\n", indent, desc[l:l+wrappedWidth])
			l += wrappedWidth
		}
		if n != len(names)-1 {
			fmt.Fprintln(w)
		}
	}
}

func PrintUsageWithWidth(width int) {
	FprintUsageWithWidth(os.Stdout, width)
}

func FprintUsage(w io.Writer) {
	FprintUsageWithWidth(w, 80)
}

func PrintUsage() {
	PrintUsageWithWidth(80)
}
