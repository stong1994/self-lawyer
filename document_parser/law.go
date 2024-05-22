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
		fmt.Printf("the %dth law %s\n", i, law.Title)
		for j, content := range law.Content {
			fmt.Printf("\t%d: %s\n", j, content)
		}
		fmt.Println("=======================================")
	}
}

type Law struct {
	Title   string
	Content []string
}

func Parse() (Laws, error) {
	data, err := os.ReadFile("./document_parser/laodongfa.md")
	if err != nil {
		return nil, err
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(data))

	var lastLaw *Law

	var laws []Law
	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node.Kind() {
		case ast.KindHeading:
			node.RemoveAttributes()
			if lastLaw != nil && len(lastLaw.Content) > 0 {
				laws = append(laws, *lastLaw)
			}
			line := node.Lines().At(0)

			lastLaw = &Law{
				Title: string(line.Value(data)),
			}
		case ast.KindParagraph:
			node.RemoveAttributes()
			line := node.Lines().At(0)
			if lastLaw != nil {
				// TODO: why goldmark got same content?
				if len(lastLaw.Content) == 0 || lastLaw.Content[len(lastLaw.Content)-1] != string(line.Value(data)) {
					lastLaw.Content = append(lastLaw.Content, string(line.Value(data)))
				}
			} else {
				fmt.Printf("No title for this content: %s\n", string(line.Value(data)))
			}
		default:
			// fmt.Printf("default: %s\n", node.Kind())
			// fmt.Println(string(node.Text(doc.Text(data))))
		}
		return ast.WalkContinue, nil
	})
	return laws, nil
}
