package dials

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Dial struct {
	Name         string
	Shorthand    string
	ValueType    string
	DefaultValue any
}

var dials map[string]*Dial
var shorthandNames map[string]string
var values map[string]any
var positionals []string

var configFiles []string

func assertValid(dial *Dial) error {
	re, err := regexp.Compile(`^([a-z]|[0-9]|-)+$`)
	if err != nil {
		log.Fatal("Error compiling regex: ", err)
	}
	if !re.MatchString(dial.Name) {
		return fmt.Errorf("Dial name '%s' is invalid! Must be lowercase a-z with dashes only.", dial.Name)
	}

	re, err = regexp.Compile(`^([a-z]|[A-Z]|[0-9])?$`)
	if err != nil {
		log.Fatal("Error compiling regex: ", err)
	}
	if !re.MatchString(dial.Shorthand) {
		return fmt.Errorf("Dial shorthand '%s' is invalid! Must be a single letter.", dial.Shorthand)
	}

	if dial.ValueType != "bool" &&
		dial.ValueType != "float" &&
		dial.ValueType != "int" &&
		dial.ValueType != "string" &&
		dial.ValueType != "floats" &&
		dial.ValueType != "ints" &&
		dial.ValueType != "strings" {
		return fmt.Errorf("Dial value type '%s' is invald! Must be one of: bool float int string floats ints strings.", dial.ValueType)
	}

	if dial.ValueType == "bool" && dial.DefaultValue != nil {
		log.Printf("Warning: You set a default value for the bool dial '%s'. It will be ignored, since bool dials always have a default value of false.\n", dial.Name)
	}
	if dial.ValueType != "bool" && dial.DefaultValue == nil {
		return fmt.Errorf("Non-bool dial '%s' must have a default value.", dial.Name)
	}

	return nil
}

func Add(dial *Dial) error {
	if dials == nil {
		dials = make(map[string]*Dial)
	}
	if shorthandNames == nil {
		shorthandNames = make(map[string]string)
	}
	if values == nil {
		values = make(map[string]any)
	}

	if err := assertValid(dial); err != nil {
		return err
	}
	if _, exists := dials[dial.Name]; exists {
		return fmt.Errorf("Dial with name '%s' already exists!", dial.Name)
	}
	if dial.Shorthand != "" {
		if _, exists := shorthandNames[dial.Shorthand]; exists {
			return fmt.Errorf("Dial with shorthand '%s' already exists!", dial.Shorthand)
		}
	}

	dials[dial.Name] = dial
	if dial.Shorthand != "" {
		shorthandNames[dial.Shorthand] = dial.Name
	}
	if dial.ValueType == "bool" {
		if err := setValue(dial, false); err != nil {
			return err
		}
	} else {
		if err := setValue(dial, dial.DefaultValue); err != nil {
			return err
		}
	}

	return nil
}

func setValue(dial *Dial, anyValue any) error {
	switch dial.ValueType {
	case "bool":
		value, ok := anyValue.(bool)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type bool!", anyValue)
		}
		values[dial.Name] = value
	case "float":
		value, ok := anyValue.(float64)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type float64!", anyValue)
		}
		values[dial.Name] = value
	case "int":
		value, ok := anyValue.(int)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type int!", anyValue)
		}
		values[dial.Name] = value
	case "string":
		value, ok := anyValue.(string)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type string!", anyValue)
		}
		values[dial.Name] = value
	case "floats":
		value, ok := anyValue.([]float64)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type []float64!", anyValue)
		}
		values[dial.Name] = value
	case "ints":
		value, ok := anyValue.([]int)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type []int!", anyValue)
		}
		values[dial.Name] = value
	case "strings":
		value, ok := anyValue.([]string)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type []string!", anyValue)
		}
		values[dial.Name] = value
	}

	return nil
}

func setStringValue(name string, str string) error {
	dial, exists := dials[name]
	if !exists {
		log.Fatalf("Setting value for dial '%s', but dial doesn't exist!", name)
	}

	switch dial.ValueType {
	case "bool":
		log.Fatal("Should never need to parse a value for a bool flag!")
	case "float":
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be a float!", dial.Name)
		}
		values[dial.Name] = value
	case "floats":
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be a float!", dial.Name)
		}
		values[dial.Name] = append(values[dial.Name].([]float64), value)
	case "int":
		value64, err := strconv.ParseInt(str, 10, 32)
		value := int(value64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be an int!", dial.Name)
		}
		values[dial.Name] = value
	case "ints":
		value64, err := strconv.ParseInt(str, 10, 32)
		value := int(value64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be an int!", dial.Name)
		}
		values[dial.Name] = append(values[dial.Name].([]int), value)
	case "string":
		values[dial.Name] = str
	case "strings":
		values[dial.Name] = append(values[dial.Name].([]string), str)
	}

	return nil
}

func parseArgs(args ...string) error {
	if len(args) <= 1 {
		return nil
	}

	needValueForName := ""
	for _, arg := range args[1:] {
		if needValueForName != "" {
			if err := setStringValue(needValueForName, arg); err != nil {
				return err
			}
			needValueForName = ""
		} else if strings.HasPrefix(arg, "--") {
			name := arg[2:]

			dial, exists := dials[name]
			if !exists {
				return fmt.Errorf("Unknown flag: --%s", name)
			}

			if dial.ValueType == "bool" {
				values[name] = true
			} else {
				needValueForName = name
			}
		} else if strings.HasPrefix(arg, "-") {
			shorthands := arg[1:]
			for c, char := range shorthands {
				shorthand := string(char)

				name, exists := shorthandNames[shorthand]
				if !exists {
					return fmt.Errorf("Unknown flag: -%s", shorthand)
				}

				dial, exists := dials[name]
				if !exists {
					log.Fatalf("Shorthand '%s' exists, but dial '%s' does not!", shorthand, name)
				}

				if dial.ValueType == "bool" {
					values[name] = true
				} else if c != len(shorthands)-1 {
					return fmt.Errorf("Invalid flag: -%s is unable to set -%s", shorthands, shorthand)
				} else {
					needValueForName = name
				}
			}
		} else {
			positionals = append(positionals, arg)
		}
	}

	return nil
}

func parseJson(data []byte) error {
	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}

	for name, anyValue := range config {
		dial, exists := dials[name]
		if !exists {
			return fmt.Errorf("Invalid JSON: Option '%s' does not exist.", name)
		}
		if err := setValue(dial, anyValue); err != nil {
			return err
		}
	}

	return nil
}

func getValue[T bool | float64 | int | string | []float64 | []int | []string](name string, valueType string) T {
	dial, exists := dials[name]
	if !exists {
		log.Fatalf("Dial '%s' does not exist!", name)
	}
	if dial.ValueType != valueType {
		log.Fatalf("Dial '%s' is not of type %s!", name, valueType)
	}
	if values[name] == nil {
		log.Fatalf("Value for '%s' not found!", name)
	}
	return values[name].(T)
}

func BoolValue(name string) bool {
	return getValue[bool](name, "bool")
}

func FloatValue(name string) float64 {
	return getValue[float64](name, "float")
}

func IntValue(name string) int {
	return getValue[int](name, "int")
}

func StringValue(name string) string {
	return getValue[string](name, "string")
}

func FloatValues(name string) []float64 {
	return getValue[[]float64](name, "floats")
}

func IntValues(name string) []int {
	return getValue[[]int](name, "ints")
}

func StringValues(name string) []string {
	return getValue[[]string](name, "strings")
}

func AddConfigFile(path string) {
	if configFiles == nil {
		configFiles = []string{path}
	} else {
		configFiles = append(configFiles, path)
	}
}

func AddHomeConfigFile(path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Warning: Could not determine location of user home directory")
	}

	AddConfigFile(home + "/" + path)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func toEnvVar(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), "-", "_")
}

func load(args ...string) {
	// 1. Config files
	for _, file := range configFiles {
		if fileExists(file) {
			file, err := os.Open(file)
			if err != nil {
				log.Printf("Failed to open file: %s\n", err)
			}
			defer file.Close()

			data, _ := io.ReadAll(file)
			if err != nil {
				log.Printf("Failed to read file: %s\n", err)
			}

			if err := parseJson(data); err != nil {
				log.Fatal(err)
			}
		}
	}

	// 2. Environment variables
	for _, dial := range dials {
		if dial.ValueType == "floats" ||
			dial.ValueType == "ints" ||
			dial.ValueType == "strings" {
			// TODO: support comma-separated environment variables
			continue
		}
		envVar := toEnvVar(dial.Name)
		value, exists := os.LookupEnv(envVar)
		if exists {
			if dial.ValueType == "bool" {
				values[dial.Name] = true
			} else if err := setStringValue(dial.Name, value); err != nil {
				log.Fatal(err)
			}
		}
	}

	// 3. Args
	if err := parseArgs(args...); err != nil {
		log.Fatal(err)
	}
}

func Load() {
	load(os.Args...)
}

func Positionals() []string {
	return positionals
}
