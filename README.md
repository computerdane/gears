# gears 

⚙️ Dead-simple configuration for Go CLI tools ⚙️

# How it Works

Define a flag that your program accepts:

```go
  gears.Add(&gears.Flag{
  	Name:         "hello-name",
  	ValueType:    "string",
  	DefaultValue: "world",
  	Shorthand:    "n",
  })
```

Add one or more locations of config files:

```go
  // Read from ~/.config/hello-world/config.json
  gears.AddHomeConfigFile(".config/hello-world/config.json")

  // Allow user to specify their own config file location
  if configFile, exists := os.LookupEnv("CONFIG_FILE"); exists {
  	gears.AddConfigFile(configFile)
  }
```

Load flags:

```go
  // Load from config files, environment, and program args
  gears.Load()

  // Collect positional arguments (non-flags)
  args := gears.Positionals()
```

Get values:

```go
  name := gears.StringValue("hello-name")

  fmt.Printf("Hello, %s!\n", name)
```

# Using Flags

Flags are always processed in the following order:

1. Configuration files (processed in the order they are added)
2. Environment variables
3. Command-line arguments

The following configuration methods all produce the same result:

1. Configuration files

```json
  {
    "hello-name": "gears"
  }
```

2. Environment variables

```sh
  HELLO_NAME="gears" hello-world
```

3. Command-line arguments

```sh
  hello-world --hello-name "gears"
  # or, using the shorthand
  hello-world -n "gears"
```

Flags processed later-on in the cycle take precedence, so command-line arguments will override environment variables, which will override config files:

```sh
  HELLO_NAME="environment" hello-world --hello-name "args"
  # Outputs: "Hello, args!"
```

# Types of Flags

There are 7 flag types: `bool` `float` `int` `string` `floats` `ints` `strings`. Set a flag's `ValueType` to select one.

The Go types for each are as follows:

| `Flag.ValueType` | Type      |
|------------------|-----------|
| bool             | bool      |
| float            | float64   |
| int              | int       |
| string           | string    |
| floats           | []float64 |
| ints             | []int     |
| strings          | []string  |

All flag types, except `bool`, must have default values. This is to ensure that when you get values, they will never be `nil`.

All `bool` flags default to `false` to ensure that boolean flags are always used as "on" switches, following a consistent pattern to avoid confusion.

To get values for each type of flag:

```go
  // Single values
  myBool := gears.BoolValue("my-bool")
  myFloat := gears.FloatValue("my-float")
  myInt := gears.IntValue("my-int")
  myString := gears.StringValue("my-string")

  // Arrays
  myFloatArray := gears.FloatValues("my-float-array")
  myIntArray := gears.IntValues("my-int-array")
  myStringArray := gears.StringValues("my-string-array")
```
