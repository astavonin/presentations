package main

import (
	"reflect"
)

type Foo struct {
	Val int
}

func (f Foo) Boo() {

}

func main() {
	f := Foo{}
	f.Boo()

	fT := reflect.TypeOf(Foo{})
	fV := reflect.New(fT)
	m := fV.MethodByName("Boo")
	m.Call(nil)
}
