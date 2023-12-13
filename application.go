package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/forPelevin/gomoji"
	"github.com/kyokomi/emoji"
)

func main() {
	http.HandleFunc("/", emojiHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/images/favicon.ico")
}

func emojiHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	emojis := gomoji.FindAll(path)
	if len(emojis) > 0 {
		for _, e := range emojis {
			emoji.Fprintln(w, e)
		}

	} else {
		t, _ := template.ParseFiles("static/notfound.gtpl")
		t.Execute(w, nil)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Return search page
		t, _ := template.ParseFiles("static/search.gtpl")
		t.Execute(w, nil)
	} else {
		// Get results of search and publish
		r.ParseForm()
		search := strings.Join(r.Form["search"], "")
		found := false
		for _, e := range gomoji.AllEmojis() {
			if strings.Contains(e.Slug, strings.ToLower(search)) {
				found = true
				emoji.Fprintln(w, e)
			}

		}
		if !found {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "No emojis found. Search again <a href=\"/search\">here</a>.")
		}
	}
}
