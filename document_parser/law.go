package document_parser

import (
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Laws []Law

func (l Laws) Print() {
	for i, law := range l {
		fmt.Printf("the %dth law %s\n", i, law.Chapter)
		for _, item := range law.Items {
			fmt.Printf("\t%s\n", item.Content)
		}
	}
}

type Item struct {
	Title   string
	Content string
}

type Law struct {
	Chapter string
	Items   []Item
}

var sourceFile = "./document_parser/laodongfa.md"

func init() {
	if s := os.Getenv("SOURCE_FILE"); s != "" {
		sourceFile = s
	}
}

func Parse() (Laws, error) {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return nil, err
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(data))

	var lastLaw *Law

	var laws []Law
	err = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node.Kind() {
		case ast.KindHeading:
			if entering {
				if lastLaw != nil {
					laws = append(laws, *lastLaw)
				}
				line := node.Lines().At(0)
				lastLaw = &Law{
					Chapter: string(line.Value(data)),
				}
			}
		case ast.KindEmphasis:
			if entering {
				title := string(node.Text(data))
				if lastLaw == nil || lastLaw.Chapter == "" {
					// log.Printf("No chapter for this content: %s\n", title)
				} else {
					lastLaw.Items = append(lastLaw.Items, Item{Title: title})
				}
			}
		case ast.KindText:
			if entering {
				// TODO: the content of KindText contains the content of KindEmphasis
				content := string(node.Text(data))
				if lastLaw == nil || len(lastLaw.Items) == 0 {
					// log.Printf("No title for this content: %s\n", content)
				} else {
					lastLaw.Items[len(lastLaw.Items)-1].Content += " " + content
				}
			}
		default:
			// fmt.Printf("default: %s\n", node.Kind())
			// fmt.Println(string(node.Text(doc.Text(data))))
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	return laws, nil
}
