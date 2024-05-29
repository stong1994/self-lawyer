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
	maxKind := ""
	maxChapter := ""
	maxItem := ""
	for i, law := range laws {
		fmt.Printf("the %d laws\n", i)
		law.Print()
		fmt.Println("-----------------------------------")

		if len(law.Kind) > len(maxKind) {
			maxKind = law.Kind
		}
		for _, chapter := range law.Chapters {
			if len(chapter.Chapter) > len(maxChapter) {
				maxChapter = chapter.Chapter
			}
			for _, item := range chapter.Items {
				if len(item.Content) > len(maxItem) {
					maxItem = item.Content
				}
			}
		}
	}
	fmt.Println("maxLenKind:", len(maxKind), maxKind)
	fmt.Println("maxLenChapter:", len(maxChapter), maxChapter)
	fmt.Println("maxLenItem:", len(maxItem), maxItem)
}
