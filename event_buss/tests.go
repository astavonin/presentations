package main

import (
	"fmt"
	"reflect"
)

type MyInt int

type Foo struct {
}

func (f Foo) Boo() {
	fmt.Println("Hi here!")
}

func kindVsType() {
	var x MyInt = 42
	var y int = 43

	y = int(x)

	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	fmt.Println(" type:", v1.Type(), "kind:", v1.Kind())
	fmt.Println(" type:", v2.Type(), "kind:", v2.Kind())

	foo := Foo{}
	v3 := reflect.ValueOf(foo)
	fmt.Println("type:", v3.Type(), "kind:", v3.Kind())
}

func rtypeAddrTest() {
	var x float64 = 3.4
	var y float64 = 3.4
	t1 := reflect.TypeOf(x)
	t2 := reflect.TypeOf(y)
	fmt.Println("type:", t1)
	fmt.Println("addr:", &t1)
	fmt.Println("addr:", &t2)
}

func reflectInfoExtraction() {
	var x float64 = 3.4
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	//fmt.Println("value:", v.Int()) // panic: reflect: call of reflect.Value.Int on float64 Value
	fmt.Println("value:", v.Float())
}

func interfaceArg(i interface{}) {
	val := i.(int)
	fmt.Println(val)
	//v := reflect.ValueOf(i)
	//fmt.Println("type:", reflect.TypeOf(i))
}

func valueToInterface() {
	var x = 3.14
	var y interface{} = x
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type(), "kind:", v.Kind())
	yType := reflect.TypeOf(y)
	fmt.Println(yType)
	z := y.(float64)
	fmt.Println(z)
	interfaceArg(y)
}

func funcToCall(data string) {
	fmt.Println("From funcToCall:", data)
}

func testReflectCall() {
	v := reflect.ValueOf(funcToCall)
	var args []reflect.Value
	args = append(args, reflect.ValueOf("test string"))
	v.Call(args)

	fmt.Println("type:", v.Type(), "kind:", v.Kind())

	args = append(args, reflect.ValueOf(42))
	//v.Call(args) // panic: reflect: Call with too many input arguments
}

func newObject() {
	f := Foo{}
	f.Boo()

	// With reflection
	fT := reflect.TypeOf(Foo{})
	fV := reflect.New(fT)
	m := fV.MethodByName("Boo")
	m.Call(nil)
}

func main() {
	//kindVsType()
	//rtypeAddrTest()
	//reflectInfoExtraction()
	//valueToInterface()
	testReflectCall()
	//var x float64 = 3.14
	//var y interface{} = x
	//z := y.(int)
	//interfaceArg(x)
}
