package main

import (
	"fmt"
	"self-lawyer/document_parser"
)

func main() {
	laws, err := document_parser.ParseAll()
	if err != nil {
		panic(err)
	}
	_ = laws
	fmt.Println("==================================================")
	for i, law := range laws {
		fmt.Printf("the %d laws\n", i)
		law.Print()
		fmt.Println("-----------------------------------")
	}
}
