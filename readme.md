## Golang | A simple configmanager program for assigning values to variables from a custom config file 


This program `configmanager` reads a configuration file (*.conf) and assigns values to a struct defined by the user. This comes in handy when you want to assign some variables in golang program from a custom configuration file. The process is simple: import `configmanager`, create your struct with **exported** fields, create a `*.conf` file or a `default.conf` file in PWD, define the variables in the config file, pass the struct reference to `AssignConfiguration` method in `configmanager`; you are good to go.

### Configuration File Format

- Comments shall start with `#`

  ```bash
  # this is a comment
  Val = 5 # this is a comment
  ```

- All variables in the configuration file must start with `Caps Letter`

- Variable must be defined with a `field name` and a `value` separated with `=`

- If a variable is defined multiple times in a single configuration file, the bottom-most value of the corresponding variable is assigned

- For variable of type slice, the values in configuration files shall be separate with "|". In case you need to change the separator character,then specify this in your program:

```go
configmanager.DefaultSliceParametersSeparator =  character/for/slice/parameter/separator // Never use "#"
```

### Struct Definition

Each field of the user-defined struct shall be exported (First letter must be in caps), so that the same can be accessed by `configmanager` program to assign values to struct fields from config file.

Following types are supported for struct fields:

> - int64
> - float64
> - string
> - bool 
> - []int64
> - []float64
> - []string

```go
type config struct {
	Val1 int64
	Val2 string
	Val3 bool
	Val4 float64
	Val5 []string
	Val6 []int64
	Val7 []float64
}
```

### configmanager usage

let's first create a file `main.go` .  Import `configmanager` from GitHub. We define a struct config with four fields: `Val1, Val2, Val3, Val4`. 

```go
package main

import (
	"fmt"
	"github.com/rrrcode9-golang-util/configmanager"
)

type config struct {
	Val1 int64
	Val2 string
	Val3 bool
	Val4 float64
	Val5 []string
	Val6 []int64
	Val7 []float64
}

func main() {
	c := config{}
	configmanager.AssignConfiguration(&c)

	fmt.Printf("%+v\n", c)

}
```

`AssignConfiguration` method takes the reference of struct variable, and then the fields will be assigned with the values from the config file.

Now create a config file `default.conf` in the present directory. 

>Note: Alternatively, `*.conf` file can be created anywhere in the system and the same can be specified while running your program passing arguments as mentioned in `Notes` below.

```conf
# file: ./default.conf
Val1 = 40
Val2 = hello
Val3 = true
Val4 = 80.7
Val5 = hi | hello | hola
Val6 = 2 | 3 | 4
Val7 = 2.5 | 5 |9.8
```

Now, run the following command:
```bash
$ go run main.go 

```

Output: `{Val1:40 Val2:hello Val3:true Val4:80.7 Val5:[hi hello hola] Val6:[2 3 4] Val7:[2.5 5 9.8]}`

### Notes: 

> Alternatively, the config file location can be specified using arguments with "--config-file-path" or "-f" and with path of config file.

```bash
$ go run main.go --config-file-path <path of config file>
or
$ go run main.go -f <path of config file>
```

> In case no config-file is specified, the `configmanager` program will search for `default.conf` file in the present working directory.


The default path for config file is `./default.conf`. In case you need to change the path (or name) of default config file, then please sepecify this in your program file: 

```go
configmanager.DefaultConfigFilePath =  path/for/default/config/file 
```
