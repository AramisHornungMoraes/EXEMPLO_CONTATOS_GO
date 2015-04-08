// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strconv"
	"fmt"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {

// hardcoded ate o fim

	var pos_livre int

	// testa para ver uma posicao livre
	for i := 1; i < 500; i++ {
		contact_file_name := strconv.Itoa(i)+".txt"
		_, err := ioutil.ReadFile(contact_file_name)
		if err != nil { pos_livre = i
		break
		}
	}
	
	fmt.Fprintf(w, "<p>[<a href=\"/create/%d\">Create</a>]</p>",pos_livre)	
	
	for i := 1; i < 500; i++ {
		contact_file_name := strconv.Itoa(i)+".txt"
		contact_content, err := ioutil.ReadFile(contact_file_name)
		if err == nil {
			
			fmt.Printf("%d - %s \n", i, contact_content)
			fmt.Fprintf(w, "<p>%d. %s - [<a href=\"/edit/%s\">edit</a>][<a href=\"/delete/%s\">delete</a>]</p>",i,contact_content, strconv.Itoa(i),strconv.Itoa(i))
		}
		if err != nil {
			pos_livre = i
		}
	}
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	//renderTemplate(w, "view", p) // nao mostra html da forma certa, somente texto pleno
}

func createHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "create", p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	os.Remove(title+".txt")
	http.Redirect(w, r, "/view/home", http.StatusFound)
}


func saveCreateHandler(w http.ResponseWriter, r *http.Request, title string) {
	
	
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	
	
	ioutil.WriteFile(title+".txt", p.Body, 0600)

	http.Redirect(w, r, "/view/home", http.StatusFound)
}


func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "create.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|create|delete|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/create/", makeHandler(createHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))
	
	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8081", nil)
}