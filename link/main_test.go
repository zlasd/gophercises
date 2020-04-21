package main

import "testing"

var testFiles = []string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}
var hrefAns = [][]string{
	{"/other-page"},
	{"https://www.twitter.com/joncalhoun", "https://github.com/gophercises"},
	{"#", "/lost", "https://twitter.com/marcusolsson"},
	{"/dog-cat"},
}

func TestMain(t *testing.T) {
	for i, fname := range testFiles {
		aList, err := parseHTML(fname)
		if err != nil {
			t.Fatal(err, fname)
		}
		for j, link := range aList {
			if link.Href != hrefAns[i][j] {
				t.Fatalf("Wrong answer in %v", fname)
			}
		}
	}
}
