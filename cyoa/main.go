package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// Forward : link
type Forward struct {
	Arc  string
	Text string
}

// Page : a story
type Page struct {
	Title   string
	Story   []string
	Options []Forward
}

func main() {
	file, _ := os.Open("gopher.json")
	storyContent := make([]byte, 1024*16)
	len, _ := file.Read(storyContent)

	story := make(map[string]Page)
	err := json.Unmarshal(storyContent[:len], &story)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(story)

	handler := adventureHandler(story)

	//http.HandleFunc("/", adventureHandler)

	fmt.Println("Listen on :8000")
	http.ListenAndServe("localhost:8000", handler)
}

func adventureHandler(story map[string]Page) http.HandlerFunc {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<div>{{range .Story}}<p>{{ . }}</p>{{else}}<p><strong>no rows</strong></p>{{end}}</div>
		<div><ul>{{range .Options}}<li><a href="/{{.Arc}}">{{.Text}}</a></li>{{end}}</ul></div>
	</body>
</html>`
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arc := r.URL.Path
		if arc == "/" {
			arc = "/intro"
		}
		arc = strings.Trim(arc, "/")
		page := story[arc]

		t.Execute(w, page)
	})
}
