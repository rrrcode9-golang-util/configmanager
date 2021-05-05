This program `configmanager` reads a configuration file (*.conf) and assign values to struct defined by the user.

### Configuration File Format

- Comments shall start wiyh `#`

  ```bash
  # this is a comment
  Val = 5 # this is a comment
  ```

- All variables in the configuration file must start with `Caps Letter`

- Variable must be defined with a `field name` and a `value` sepearated with `=`

- If a variable is defined multiple times in a single configuration file, the bottom-most value of the corresponding variable is assigned

### Struct Definition

Each field of the user defined struct shall be exported (First letter must be in caps), so that the same can be accesed by `configmanager` program to assign values to struct fields from config file.

Following types are supported for struct fields:

> - int64
> - float64
> - string
> - bool 

```go
type config struct {
	Val1 int64
	Val2 string
	Val3 bool
	Val4 float64
}
```

### `configmanager` usage

let's first create a file `main.go` .  We define a struct config with four fields: Val1, Val2, Val3, Val4. The config file location can be specified using arguments with "--config-file-path" or "-f" and with path of config file.

```bash
$ go run main.go --config-file-path <path to configfile>
or
$ go run main.go -f <path to configfile>
```

In case no config-file is specified, the `configmanager` program will search for `default.conf` file in the present working directory.

Example:

```go
package main

import (
	"fmt"
	"github.com/rrrcode9/configmanager"
)

type config struct {
	Val1 int64
	Val2 string
	Val3 bool
	val4 float64
}

func main() {
	c := config{}
	configmanager.AssignConfiguration(&c)

	fmt.Println(c)

}
```

`AssignConfiguration` method takes the reference of struct variable, and then the fields will be assigned with the values from the config file.

To specify config file path in environment variable, assign config file path to CONFIG_FILE_PATH variable
