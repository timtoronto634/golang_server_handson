package main

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// // fmt.Fprintf(w, "Hello, %q", r.URL.Path[2:])
	title := r.URL.Path[6:]
	p, _ := loadPage(title)
	// fmt.Fprintf(w, "<h1>%s</h1><p>[<a href=\"/edit/%s\">edit</a>]</p><div>%s</div>", p.Title, p.Title, p.Body)

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[6:]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[6:]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
