package gears

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

type Flag struct {
	Name         string
	Shorthand    string
	ValueType    string
	DefaultValue any
}

var flags map[string]*Flag
var shorthandNames map[string]string
var values map[string]any
var positionals []string

var configFiles []string

func assertValid(flag *Flag) error {
	re, err := regexp.Compile(`^([a-z]|[0-9]|-)+$`)
	if err != nil {
		log.Fatal("Error compiling regex: ", err)
	}
	if !re.MatchString(flag.Name) {
		return fmt.Errorf("Flag name '%s' is invalid! Must be lowercase a-z with dashes only.", flag.Name)
	}

	re, err = regexp.Compile(`^([a-z]|[A-Z]|[0-9])?$`)
	if err != nil {
		log.Fatal("Error compiling regex: ", err)
	}
	if !re.MatchString(flag.Shorthand) {
		return fmt.Errorf("Flag shorthand '%s' is invalid! Must be a single letter.", flag.Shorthand)
	}

	if flag.ValueType != "bool" &&
		flag.ValueType != "float" &&
		flag.ValueType != "int" &&
		flag.ValueType != "string" &&
		flag.ValueType != "floats" &&
		flag.ValueType != "ints" &&
		flag.ValueType != "strings" {
		return fmt.Errorf("Flag value type '%s' is invald! Must be one of: bool float int string floats ints strings.", flag.ValueType)
	}

	if flag.ValueType == "bool" && flag.DefaultValue != nil {
		log.Printf("Warning: You set a default value for the bool flag '%s'. It will be ignored, since bool flags always have a default value of false.\n", flag.Name)
	}
	if flag.ValueType != "bool" && flag.DefaultValue == nil {
		return fmt.Errorf("Non-bool flag '%s' must have a default value.", flag.Name)
	}

	return nil
}

func Add(flag *Flag) error {
	if flags == nil {
		flags = make(map[string]*Flag)
	}
	if shorthandNames == nil {
		shorthandNames = make(map[string]string)
	}
	if values == nil {
		values = make(map[string]any)
	}

	if err := assertValid(flag); err != nil {
		return err
	}
	if _, exists := flags[flag.Name]; exists {
		return fmt.Errorf("Flag with name '%s' already exists!", flag.Name)
	}
	if flag.Shorthand != "" {
		if _, exists := shorthandNames[flag.Shorthand]; exists {
			return fmt.Errorf("Flag with shorthand '%s' already exists!", flag.Shorthand)
		}
	}

	flags[flag.Name] = flag
	if flag.Shorthand != "" {
		shorthandNames[flag.Shorthand] = flag.Name
	}
	if flag.ValueType == "bool" {
		if err := setValue(flag, false); err != nil {
			return err
		}
	} else {
		if err := setValue(flag, flag.DefaultValue); err != nil {
			return err
		}
	}

	return nil
}

func setValue(flag *Flag, anyValue any) error {
	switch flag.ValueType {
	case "bool":
		value, ok := anyValue.(bool)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type bool!", anyValue)
		}
		values[flag.Name] = value
	case "float":
		value, ok := anyValue.(float64)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type float64!", anyValue)
		}
		values[flag.Name] = value
	case "int":
		value, ok := anyValue.(int)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type int!", anyValue)
		}
		values[flag.Name] = value
	case "string":
		value, ok := anyValue.(string)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type string!", anyValue)
		}
		values[flag.Name] = value
	case "floats":
		value, ok := anyValue.([]float64)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type []float64!", anyValue)
		}
		values[flag.Name] = value
	case "ints":
		value, ok := anyValue.([]int)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type []int!", anyValue)
		}
		values[flag.Name] = value
	case "strings":
		value, ok := anyValue.([]string)
		if !ok {
			return fmt.Errorf("Value '%v' is not of type []string!", anyValue)
		}
		values[flag.Name] = value
	}

	return nil
}

func SetValue(name string, value any) {
	flag, exists := flags[name]
	if !exists {
		log.Fatalf("Flag '%s' does not exist!", name)
	}
	setValue(flag, value)
}

func setJsonValue(flag *Flag, raw json.RawMessage) error {
	switch flag.ValueType {
	case "bool":
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type bool!", string(raw))
		}
		values[flag.Name] = value
	case "float":
		var value float64
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type float64!", string(raw))
		}
		values[flag.Name] = value
	case "int":
		var value int
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type int!", string(raw))
		}
		values[flag.Name] = value
	case "string":
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type string!", string(raw))
		}
		values[flag.Name] = value
	case "floats":
		var value []float64
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type []float64!", string(raw))
		}
		values[flag.Name] = value
	case "ints":
		var value []int
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type []int!", string(raw))
		}
		values[flag.Name] = value
	case "strings":
		var value []string
		if err := json.Unmarshal(raw, &value); err != nil {
			return fmt.Errorf("JSON value '%s' is not of type []string!", string(raw))
		}
		values[flag.Name] = value
	}

	return nil
}

func setStringValue(name string, str string) error {
	flag, exists := flags[name]
	if !exists {
		log.Fatalf("Setting value for flag '%s', but flag doesn't exist!", name)
	}

	switch flag.ValueType {
	case "bool":
		log.Fatal("Should never need to parse a value for a bool flag!")
	case "float":
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be a float!", flag.Name)
		}
		values[flag.Name] = value
	case "floats":
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be a float!", flag.Name)
		}
		values[flag.Name] = append(values[flag.Name].([]float64), value)
	case "int":
		value64, err := strconv.ParseInt(str, 10, 32)
		value := int(value64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be an int!", flag.Name)
		}
		values[flag.Name] = value
	case "ints":
		value64, err := strconv.ParseInt(str, 10, 32)
		value := int(value64)
		if err != nil {
			return fmt.Errorf("Value for '%s' must be an int!", flag.Name)
		}
		values[flag.Name] = append(values[flag.Name].([]int), value)
	case "string":
		values[flag.Name] = str
	case "strings":
		values[flag.Name] = append(values[flag.Name].([]string), str)
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

			flag, exists := flags[name]
			if !exists {
				return fmt.Errorf("Unknown flag: --%s", name)
			}

			if flag.ValueType == "bool" {
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

				flag, exists := flags[name]
				if !exists {
					log.Fatalf("Shorthand '%s' exists, but flag '%s' does not!", shorthand, name)
				}

				if flag.ValueType == "bool" {
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
	var config map[string]json.RawMessage
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}

	for name, raw := range config {
		flag, exists := flags[name]
		if !exists {
			return fmt.Errorf("Invalid JSON: Option '%s' does not exist.", name)
		}
		if err := setJsonValue(flag, raw); err != nil {
			return err
		}
	}

	return nil
}

func getValue[T bool | float64 | int | string | []float64 | []int | []string](name string, valueType string) T {
	flag, exists := flags[name]
	if !exists {
		log.Fatalf("Flag '%s' does not exist!", name)
	}
	if flag.ValueType != valueType {
		log.Fatalf("Flag '%s' is not of type %s!", name, valueType)
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
	for _, flag := range flags {
		if flag.ValueType == "floats" ||
			flag.ValueType == "ints" ||
			flag.ValueType == "strings" {
			// TODO: support comma-separated environment variables
			continue
		}
		envVar := toEnvVar(flag.Name)
		value, exists := os.LookupEnv(envVar)
		if exists {
			if flag.ValueType == "bool" {
				values[flag.Name] = true
			} else if err := setStringValue(flag.Name, value); err != nil {
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
