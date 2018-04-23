package main

import (
       "html/template"
       "io/ioutil"
       "log"
       "net/http"
       "regexp"
       "errors"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


type Page struct {
     Title string
     Body []byte
}

func (p *Page) save() error {
     filename := p.Title + ".txt"
     return ioutil.WriteFile(filename, p.Body, 0600)
}

/*func loadPage(title string) *Page {
     filename := title + ".txt"
     body, _ := ioutil.ReadFile(filename)
     return &Page{Title: title, Body: body}
}*/

func loadPage(title string) (*Page, error) {
     filename := title + ".txt"
     body, err := ioutil.ReadFile(filename)
     if err != nil {
     	return nil, err
     }
     return &Page{Title: title, Body: body}, nil
}

/*func main() {
     p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
     p1.save()
     p2, _ := loadPage("TestPage")
     fmt.Println(string(p2.Body))
}*/

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
     // title := r.URL.Path[len("/view/"):]
     // title, err := getTitle(w, r)
     // if err != nil {
     // 	return
     // }
     p, err := loadPage(title)
     if err != nil {
     	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	return
     }
     renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
     // title := r.URL.Path[len("/edit/"):]
     // title, err := getTitle(w, r)
     // if err != nil {
     // 	return
     // }
     p, err := loadPage(title)
     if err != nil {
     	p = &Page{Title: title}
     }
     renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
//     title := r.URL.Path[len("/save/"):]
     // title, err := getTitle(w, r)
     // if err != nil {
     // 	return
     // }
     body := r.FormValue("body")
     p := &Page{Title: title, Body: []byte(body)}
     err := p.save()
     if err != nil {
     	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
     }
     http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
     return func(w http.ResponseWriter, r *http.Request) {
     	// Here we will extract the page title from the request,
   	// and call the prived handler 'fn'
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
	   http.NotFound(w, r)
	   return
	}
	fn(w, r, m[2])
     }
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
     // t, err := template.ParseFiles(tmpl + ".html")
     err := templates.ExecuteTemplate(w, tmpl+".html", p)
     if err != nil {
     	http.Error(w, err.Error(), http.StatusInternalServerError)
//	return
     }
//     err = t.Execute(w, p)
//     if err != nil {
//     	http.Error(w, err.Error(), http.StatusInternalServerError)
//     }
 }

func getTitle(w http.ResponseWriter, r *http.Request) (string, error){
     m := validPath.FindStringSubmatch(r.URL.Path)
     if m == nil {
     	http.NotFound(w, r)
	return "", errors.New("Invalid Page Title")
     }
     return m[2], nil
}

func main() {
     http.HandleFunc("/view/", makeHandler(viewHandler))
     http.HandleFunc("/edit/", makeHandler(editHandler))
     http.HandleFunc("/save/", makeHandler(saveHandler))
     
     log.Fatal(http.ListenAndServe(":8080", nil))
}