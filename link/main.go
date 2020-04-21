package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "f", "ex1.html", "test file path")
	flag.Parse()

	aList, _ := parseHTML(fileName)
	for _, a := range aList {
		fmt.Println(a)
	}
}

// Link ...
type Link struct {
	Href string
	Text string
}

func parseHTML(fileName string) ([]Link, error) {
	htmlFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// content, _ := ioutil.ReadFile(fileName)
	// fmt.Println(string(content))

	doc, err := html.Parse(htmlFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	aList := traverse(doc)
	return aList, nil
}

func traverse(n *html.Node) []Link {
	aList := make([]Link, 0)
	if n.Type == html.ElementNode && n.Data == "a" {
		aList = append(aList, parseLink(n))
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			aList = append(aList, traverse(c)...)
		}
	}
	return aList
}

func parseLink(n *html.Node) Link {
	var href string
	for _, a := range n.Attr {
		if a.Key == "href" {
			href = a.Val
		}
	}
	text := parseText(n)
	return Link{href, text}
}

func parseText(n *html.Node) string {
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		} else {
			text += parseText(c)
		}
	}
	return text
}
