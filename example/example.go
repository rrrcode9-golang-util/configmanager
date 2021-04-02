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
}

func main() {
	configs := config{}

	//specify default config file path | Optional [by default it is - 'default.conf' in PWD]
	configmanager.DefaultConfigFilePath = "./default.conf"

	// assign configuration
	configmanager.AssignConfiguration(&configs)

	fmt.Printf("%+v\n", configs)

}
