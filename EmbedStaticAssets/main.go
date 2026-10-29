package main

import (
	_ "embed"
	"fmt"
)

//go:embed hello.txt
var data string

func main() {
	fmt.Println(data)
}
