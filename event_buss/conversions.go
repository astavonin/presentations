package main

import (
	"fmt"
)

func funcWithInterface(i interface{}) {
	val := i.(int)
	fmt.Println(val)
}

func main() {
	var x float64 = 3.14
	var y interface{} = x
	z := y.(int)
	fmt.Println(z)
}
