package main

import "fmt"

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

func incr(n *int)  {
	(*n)++
}

func main() {

	var n int
	n = 5
	incr(&n)
	fmt.Println(n)
}