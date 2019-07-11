package main

import (
	"fmt"
	"reflect"
)

type Bus interface {
	Subscribe(event string, handler interface{}) error
	Unsubscribe(event string, handler interface{}) error
	Publish(event string, args ...interface{})
}

type EventBus struct {
	handlers map[string][]reflect.Value
}

func (b *EventBus) Subscribe(event string, handler interface{}) error {
	if !(reflect.TypeOf(handler).Kind() == reflect.Func) {
		return fmt.Errorf("invalid type %s, reflect.Func is expected",
			reflect.TypeOf(handler).Kind())
	}
	b.handlers[event] = append(b.handlers[event], reflect.ValueOf(handler))

	return nil
}

func (b *EventBus) handlerIdx(event string, handler reflect.Value) int {
	for idx, h := range b.handlers[event] {
		if h == handler {
			return idx
		}
	}

	return -1
}

func (b *EventBus) Unsubscribe(event string, handler interface{}) error {
	if _, ok := b.handlers[event]; !ok {
		return fmt.Errorf("event %s doesn't exist", event)
	}
	idx := b.handlerIdx(event, reflect.ValueOf(handler))
	b.handlers[event] = append(b.handlers[event][:idx], b.handlers[event][idx+1:]...)
	return nil
}

func (b *EventBus) Publish(event string, args ...interface{}) {
	argsVal := make([]reflect.Value, 0)
	for _, arg := range args {
		argsVal = append(argsVal, reflect.ValueOf(arg))
	}
	if handlers, ok := b.handlers[event]; ok {
		for _, h := range handlers {
			h.Call(argsVal)
		}
	}
}

func NewBus() Bus {
	return &EventBus{
		handlers: make(map[string][]reflect.Value),
	}
}

func main() {
	b := NewBus()

	_ = b.Subscribe("foo", func(str string, num int) {
		fmt.Printf("(1) foo(%s, %d)\n", str, num)
	})
	f2 := func(str string, num int) {
		fmt.Printf("(2) foo(%s, %d)\n", str, num)
	}
	_ = b.Subscribe("foo", f2)
	_ = b.Subscribe("boo", func(str string) {
		fmt.Printf("(3) foo(%s)\n", str)
	})

	b.Publish("foo", "string", 42)
	_ = b.Unsubscribe("foo", f2)
	b.Publish("foo", "string1", 42)
	b.Publish("boo", "string2")
}
