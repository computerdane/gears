package gears

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func reset() {
	flags = nil
	shorthandNames = nil
	values = nil
	positionals = nil
	configFiles = nil
}

func TestAssertValid(t *testing.T) {
	// Missing parameters
	if err := assertValid(&Flag{}); err == nil {
		t.Error("empty flag is valid; want invalid")
	}
	if err := assertValid(&Flag{Name: "test"}); err == nil {
		t.Error("flag with no ValueType is valid; want invalid")
	}
	if err := assertValid(&Flag{Name: "test", ValueType: "string"}); err == nil {
		t.Error("string flag with no DefaultValue is valid; want invalid")
	}

	// flag name
	if err := assertValid(&Flag{Name: "invalid name 1", ValueType: "bool"}); err == nil {
		t.Error("flag with name 'invalid name 1' is valid; want invalid")
	}
	if err := assertValid(&Flag{Name: "Invalid-name-1", ValueType: "bool"}); err == nil {
		t.Error("flag with name 'Invalid-name-1' is valid; want invalid")
	}
	if err := assertValid(&Flag{Name: "valid-name-1", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with name 'valid-name-1' is invalid; want valid")
	}
	if err := assertValid(&Flag{Name: "valid", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with name 'valid' is invalid; want valid")
	}
	if err := assertValid(&Flag{Name: "v", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with name 'v' is invalid; want valid")
	}
	if err := assertValid(&Flag{Name: "1", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with name '1' is invalid; want valid")
	}

	// flag shorthand
	if err := assertValid(&Flag{Name: "valid-name-1", Shorthand: "vv", ValueType: "bool"}); err == nil {
		t.Error("flag with shorthand 'vv' is valid; want invalid")
	}
	if err := assertValid(&Flag{Name: "valid-name-1", Shorthand: " ", ValueType: "bool"}); err == nil {
		t.Error("flag with shorthand ' ' is valid; want invalid")
	}
	if err := assertValid(&Flag{Name: "valid-name-1", Shorthand: "v", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with shorthand 'v' is invalid; want valid")
	}
	if err := assertValid(&Flag{Name: "valid-name-1", Shorthand: "V", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with shorthand 'V' is invalid; want valid")
	}
	if err := assertValid(&Flag{Name: "valid-name-1", Shorthand: "1", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("flag with shorthand '1' is invalid; want valid")
	}
}

func TestAdd(t *testing.T) {
	// Default values
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "string", DefaultValue: 1}); err == nil {
		t.Error("string flag with DefaultValue of 1 is valid; want invalid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "int", DefaultValue: 1.0}); err == nil {
		t.Error("int flag with DefaultValue of 1.0 is valid; want invalid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "float", DefaultValue: 1}); err == nil {
		t.Error("float flag with DefaultValue of 1 is valid; want invalid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "strings", DefaultValue: "value"}); err == nil {
		t.Error("strings flag with DefaultValue of 'value' is valid; want invalid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "ints", DefaultValue: 1}); err == nil {
		t.Error("ints flag with DefaultValue of 1 is valid; want invalid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "floats", DefaultValue: 1.0}); err == nil {
		t.Error("floats flag with DefaultValue of 1.0 is valid; want invalid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "bool"}); err != nil {
		fmt.Println(err)
		t.Error("bool flag with no DefaultValue is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "string", DefaultValue: ""}); err != nil {
		fmt.Println(err)
		t.Error("string flag with DefaultValue of '' is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "string", DefaultValue: "test value"}); err != nil {
		fmt.Println(err)
		t.Error("string flag with DefaultValue of 'test value' is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "int", DefaultValue: 1}); err != nil {
		fmt.Println(err)
		t.Error("int flag with DefaultValue of 1 is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "float", DefaultValue: 1.0}); err != nil {
		fmt.Println(err)
		t.Error("float flag with DefaultValue of 1.0 is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "strings", DefaultValue: []string{}}); err != nil {
		fmt.Println(err)
		t.Error("strings flag with DefaultValue of []string{} is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "ints", DefaultValue: []int{}}); err != nil {
		fmt.Println(err)
		t.Error("ints flag with DefaultValue of []int{} is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "floats", DefaultValue: []float64{}}); err != nil {
		fmt.Println(err)
		t.Error("floats with DefaultValue of []float64{} is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "strings", DefaultValue: []string{"test"}}); err != nil {
		fmt.Println(err)
		t.Error("strings flag with DefaultValue of []string{'test'} is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "ints", DefaultValue: []int{1}}); err != nil {
		fmt.Println(err)
		t.Error("ints flag with DefaultValue of []int{1} is invalid; want valid")
	}
	reset()
	if err := Add(&Flag{Name: "test", ValueType: "floats", DefaultValue: []float64{1.0}}); err != nil {
		fmt.Println(err)
		t.Error("floats with DefaultValue of []float64{1.0} is invalid; want valid")
	}
	reset()

	// Empty shorthand
	if err := Add(&Flag{Name: "alsdjkf", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "sljkdf", ValueType: "bool"}); err != nil {
		t.Error("Add two flags with empty shorthand failed; wanted success")
	}

	// Duplicates
	if err := Add(&Flag{Name: "test", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "test", ValueType: "string", DefaultValue: ""}); err == nil {
		t.Error("Add flag with duplicate name was successful; wanted failure")
	}
	if err := Add(&Flag{Name: "abc", Shorthand: "a", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "def", Shorthand: "a", ValueType: "string", DefaultValue: ""}); err == nil {
		t.Error("Add flag with duplicate shorthand was successful; wanted failure")
	}

	// Values
	if err := Add(&Flag{Name: "mybool", ValueType: "bool", DefaultValue: true}); err != nil {
		log.Fatal(err)
	}
	if values["mybool"] != false {
		t.Error("bool flag with true default value had an actual value of true; wanted false")
	}
	if err := Add(&Flag{Name: "mystring", ValueType: "string", DefaultValue: "my value"}); err != nil {
		log.Fatal(err)
	}
	if values["mystring"] != "my value" {
		t.Error("value for string flag did not match default value")
	}

	// Empty values
	if _, exists := values[""]; exists {
		t.Error("key '' exists in values; expected none")
	}
	if _, exists := shorthandNames[""]; exists {
		t.Error("key '' exists in shorthands; expected none")
	}
}

func TestSetStringValue(t *testing.T) {
	// float/floats
	if err := Add(&Flag{Name: "test-set-value-float", ValueType: "float", DefaultValue: 0.0}); err != nil {
		log.Fatal(err)
	}
	if err := setStringValue("test-set-value-float", "1.0"); err != nil {
		t.Errorf("setValue float 1.0 failed: %v", err)
	}
	if values["test-set-value-float"] != 1.0 {
		t.Error("values['test-set-value-float'] is not 1.0")
	}
	if err := Add(&Flag{Name: "test-set-value-floats", ValueType: "floats", DefaultValue: []float64{}}); err != nil {
		log.Fatal(err)
	}
	if err := setStringValue("test-set-value-floats", "1.0"); err != nil {
		t.Errorf("setValue floats 1.0 failed: %v", err)
	}
	if values["test-set-value-floats"].([]float64)[0] != 1.0 {
		t.Error("values['test-set-value-floats'][0] is not 1.0")
	}

	// int/ints
	if err := Add(&Flag{Name: "test-set-value-int", ValueType: "int", DefaultValue: 0}); err != nil {
		log.Fatal(err)
	}
	if err := setStringValue("test-set-value-int", "1"); err != nil {
		t.Errorf("setValue int 1 failed: %v", err)
	}
	if values["test-set-value-int"] != 1 {
		t.Error("values['test-set-value-int'] is not 1")
	}
	if err := Add(&Flag{Name: "test-set-value-ints", ValueType: "ints", DefaultValue: []int{}}); err != nil {
		log.Fatal(err)
	}
	if err := setStringValue("test-set-value-ints", "1"); err != nil {
		t.Errorf("setValue ints 1 failed: %v", err)
	}
	if values["test-set-value-ints"].([]int)[0] != 1 {
		t.Error("values['test-set-value-ints'][0] is not 1")
	}

	// string/strings
	if err := Add(&Flag{Name: "test-set-value-string", ValueType: "string", DefaultValue: ""}); err != nil {
		log.Fatal(err)
	}
	if err := setStringValue("test-set-value-string", "a"); err != nil {
		t.Errorf("setValue string 'a' failed: %v", err)
	}
	if values["test-set-value-string"] != "a" {
		t.Error("values['test-set-value-string'] is not 'a'")
	}
	if err := Add(&Flag{Name: "test-set-value-strings", ValueType: "strings", DefaultValue: []string{}}); err != nil {
		log.Fatal(err)
	}
	if err := setStringValue("test-set-value-strings", "a"); err != nil {
		t.Errorf("setValue strings 'a' failed: %v", err)
	}
	if values["test-set-value-strings"].([]string)[0] != "a" {
		t.Error("values['test-set-value-strings'][0] is not 'a'")
	}
}

func TestParseArgs(t *testing.T) {
	reset()

	if err := Add(&Flag{Name: "mybool-a", Shorthand: "a", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "mybool-b", Shorthand: "b", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "mystring", Shorthand: "s", ValueType: "string", DefaultValue: ""}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "myint", Shorthand: "i", ValueType: "int", DefaultValue: 0}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "myfloat", Shorthand: "f", ValueType: "float", DefaultValue: 0.0}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "mystrings", Shorthand: "S", ValueType: "strings", DefaultValue: []string{}}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "myints", Shorthand: "I", ValueType: "ints", DefaultValue: []int{}}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "myfloats", Shorthand: "F", ValueType: "floats", DefaultValue: []float64{}}); err != nil {
		log.Fatal(err)
	}

	// Big complex but successful command
	if err := parseArgs("cmd", "pos0", "--mybool-a", "--mystring", "a", "pos1", "--myint", "1", "--myfloat", "1.0", "--mystrings", "a", "pos2", "--mystrings", "b", "--myints", "1", "--myints", "2", "--myfloats", "1.0", "--myfloats", "2.0", "pos3", "pos4"); err != nil {
		log.Fatal(err)
	}
	if BoolValue("mybool-a") != true {
		t.Error("mybool-a is false; expected true")
	}
	if StringValue("mystring") != "a" {
		t.Error("mystring is not 'a'")
	}
	if IntValue("myint") != 1 {
		t.Error("myint is not 1")
	}
	if FloatValue("myfloat") != 1.0 {
		t.Error("myfloat is not 1")
	}
	stringValues := StringValues("mystrings")
	if stringValues[0] != "a" || stringValues[1] != "b" {
		t.Error("mystrings is not [a, b]")
	}
	intValues := IntValues("myints")
	if intValues[0] != 1 || intValues[1] != 2 {
		t.Error("myints is not [0, 1]")
	}
	floatValues := FloatValues("myfloats")
	if floatValues[0] != 1.0 || floatValues[1] != 2.0 {
		t.Error("myfloats is not [0.0, 1.0]")
	}
	if positionals[0] != "pos0" || positionals[1] != "pos1" || positionals[2] != "pos2" || positionals[3] != "pos3" || positionals[4] != "pos4" {
		t.Error("failed to properly find positional arguments")
	}
	values["mybool-a"] = false
	values["mystring"] = ""
	values["myint"] = 0
	values["myfloat"] = 0.0
	values["mystrings"] = []string{}
	values["myints"] = []int{}
	values["myfloats"] = []float64{}
	positionals = []string{}

	// Big complex but successful command with shorthands
	if err := parseArgs("cmd", "pos0", "-a", "-s", "a", "pos1", "-i", "1", "-f", "1.0", "-S", "a", "pos2", "-S", "b", "-I", "1", "-I", "2", "-F", "1.0", "-F", "2.0", "pos3", "pos4"); err != nil {
		log.Fatal(err)
	}
	if BoolValue("mybool-a") != true {
		t.Error("mybool-a is false; expected true")
	}
	if StringValue("mystring") != "a" {
		t.Error("mystring is not 'a'")
	}
	if IntValue("myint") != 1 {
		t.Error("myint is not 1")
	}
	if FloatValue("myfloat") != 1.0 {
		t.Error("myfloat is not 1")
	}
	stringValues = StringValues("mystrings")
	if stringValues[0] != "a" || stringValues[1] != "b" {
		t.Error("mystrings is not [a, b]")
	}
	intValues = IntValues("myints")
	if intValues[0] != 1 || intValues[1] != 2 {
		t.Error("myints is not [0, 1]")
	}
	floatValues = FloatValues("myfloats")
	if floatValues[0] != 1.0 || floatValues[1] != 2.0 {
		t.Error("myfloats is not [0.0, 1.0]")
	}
	if positionals[0] != "pos0" || positionals[1] != "pos1" || positionals[2] != "pos2" || positionals[3] != "pos3" || positionals[4] != "pos4" {
		t.Error("failed to properly find positional arguments")
	}
	values["mybool-a"] = false
	values["mystring"] = ""
	values["myint"] = 0
	values["myfloat"] = 0.0
	values["mystrings"] = []string{}
	values["myints"] = []int{}
	values["myfloats"] = []float64{}
	positionals = []string{}

	// Combining bool shorthands
	if err := parseArgs("cmd", "-ba"); err != nil {
		log.Fatal(err)
	}
	if BoolValue("mybool-a") != true {
		t.Error("mybool-a is false; expected true")
	}
	if BoolValue("mybool-b") != true {
		t.Error("mybool-b is false; expected true")
	}
	values["mybool-a"] = false
	values["mybool-b"] = false

	// Combining bool shorthands with other shorthand
	if err := parseArgs("cmd", "-bas", "a"); err != nil {
		log.Fatal(err)
	}
	if BoolValue("mybool-a") != true {
		t.Error("mybool-a is false; expected true")
	}
	if BoolValue("mybool-b") != true {
		t.Error("mybool-b is false; expected true")
	}
	if StringValue("mystring") != "a" {
		t.Error("mystring is not 'a'")
	}
	values["mybool-a"] = false
	values["mybool-b"] = false
	values["mystrings"] = []string{}

	// No flags
	if err := parseArgs("cmd", "test"); err != nil {
		t.Error("'cmd test' failed; expected success")
	}

	// Invalid commands
	if err := parseArgs("cmd", "-asb", "a"); err == nil {
		t.Error("'cmd -asb a' succeeded; expected failure")
	}
	if err := parseArgs("cmd", "--unknown"); err == nil {
		t.Error("'cmd --unknown' succeeded; expected failure")
	}
	if err := parseArgs("cmd", "-u"); err == nil {
		t.Error("'cmd -u' succeeded; expected failure")
	}
}

func TestLoad(t *testing.T) {
	reset()

	if err := Add(&Flag{Name: "my-str", Shorthand: "s", ValueType: "string", DefaultValue: ""}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "my-bool", Shorthand: "b", ValueType: "bool"}); err != nil {
		log.Fatal(err)
	}
	if err := Add(&Flag{Name: "my-floats", Shorthand: "f", ValueType: "floats", DefaultValue: []float64{}, EnvVarDelimiter: ","}); err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Faild to get cwd: ", err)
	}

	configPath := cwd + "/test.json"
	config := `{"my-str":"json","my-floats":[1.0,2.0]}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		log.Fatal("Failed to write config: ", err)
	}
	config2Path := cwd + "/test2.json"
	config2 := `{"my-str":"json2"}`
	if err := os.WriteFile(config2Path, []byte(config2), 0644); err != nil {
		log.Fatal("Failed to write config: ", err)
	}

	// Config files
	AddConfigFile(configPath)
	values = make(map[string]any)
	load()
	if StringValue("my-str") != "json" {
		t.Error("setting my-str using JSON config failed")
	}
	myFloats := FloatValues("my-floats")
	if myFloats[0] != 1.0 || myFloats[1] != 2.0 {
		t.Error("setting my-floats using JSON config failed")
	}
	AddConfigFile(config2Path)
	values = make(map[string]any)
	load()
	if StringValue("my-str") != "json2" {
		t.Error("setting my-str using JSON config2 failed")
	}

	// Environment
	values = make(map[string]any)
	os.Setenv("MY_STR", "env")
	load()
	if StringValue("my-str") != "env" {
		t.Error("setting my-str using environment failed")
	}

	values = make(map[string]any)
	os.Setenv("MY_BOOL", "1")
	load()
	if BoolValue("my-bool") != true {
		t.Error("setting my-bool using environment failed")
	}

	values["my-bool"] = false
	os.Unsetenv("MY_BOOL")
	load()
	if BoolValue("my-bool") != false {
		t.Error("unsetting my-bool using environment failed")
	}

	values = make(map[string]any)
	configFiles = nil
	os.Setenv("MY_FLOATS", "2.0,3.0")
	load()
	myFloats = FloatValues("my-floats")
	if myFloats[0] != 2.0 || myFloats[1] != 3.0 {
		t.Error("setting my-floats using JSON config failed")
	}

	// Args
	values = make(map[string]any)
	load("cmd", "-s", "arg")
	if StringValue("my-str") != "arg" {
		t.Error("setting my-str using args failed")
	}

	os.Remove(configPath)
	os.Remove(config2Path)
}
