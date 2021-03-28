package main

import (
	"github.com/rrrcode9/configmanager"
	"fmt"

)


type config struct {
	Val1 int64
	Val2 string
	Val3 bool
	Val4 float64
}
func main(){
	c := config{}
	
	configmanager.AssignConfiguration(&c)
	fmt.Println(c)

}
