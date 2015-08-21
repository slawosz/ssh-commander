package main

import (
	"encoding/json"
	"fmt"
)

var json1 = `
{ "Foo": "bar" }
`

var json2 = `
{ "Foo": "bar", "Foo2": "baz" }
`

type F struct {
	Foo string
}

type F2 struct {
	F
	F2 string
}

func main() {
	f := &F{}
	f2 := &F2{}
	json.Unmarshal([]byte(json1), f)
	json.Unmarshal([]byte(json2), f2)
	fmt.Printf("%v\n", f)
	fmt.Printf("%v\n", f2)
}
