package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	reversedString := reverse.String("Hello, OTUS!")
	fmt.Println(reversedString)
}
