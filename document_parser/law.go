package document_parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Laws struct {
	Kind     string
	Chapters []Law
}

func (l Laws) Print() {
	fmt.Println("Law Kind:", l.Kind)
	for i, law := range l.Chapters {
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

var sourceFile = "./law_docs"

func init() {
	if s := os.Getenv("SOURCE_FILE"); s != "" {
		sourceFile = s
	}
}

func ParseAll() ([]Laws, error) {
	var rst []Laws
	err := filepath.Walk(sourceFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		laws, err := ParseFile(path, info)
		if err != nil {
			return err
		}
		rst = append(rst, laws)
		return nil
	})
	return rst, err
}

func ParseFile(path string, info os.FileInfo) (Laws, error) {
	laws := Laws{
		Kind: strings.TrimSuffix(info.Name(), ".md"),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return laws, err
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(data))

	var lastLaw *Law

	var chapters []Law
	err = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node.Kind() {
		case ast.KindHeading:
			if entering {
				if lastLaw != nil {
					chapters = append(chapters, *lastLaw)
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
		return laws, err
	}

	laws.Chapters = chapters
	return laws, nil
}
