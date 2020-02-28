package main

import (
	"encoding/json"
	"fmt"
)

var (
	mp map[string]*class
)


type class struct {
	A int    `json:"a"`
	B string `json:"b"`
	//c tl.Attr
}

type response1 struct {
	Page   int
	Fruits []string
}

func main() {
	cl := &class{
		A: 123,
		B: "asdf",
	}
	asdf, _ := json.Marshal(cl)
	fmt.Println(string(asdf))


	res1D := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))
}