package gears

import "fmt"

func FishCompletions(command string) string {
	completions := ""
	for _, flag := range flags {
		completion := fmt.Sprintf(`complete -c "%s" -l "%s"`, command, flag.Name)
		if flag.Shorthand != "" {
			completion += fmt.Sprintf(` -s "%s"`, flag.Shorthand)
		}
		if flag.Description != "" {
			completion += fmt.Sprintf(` -d "%s"`, flag.Description)
		}
		completions += completion + "\n"
	}
	return completions
}
