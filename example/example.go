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
