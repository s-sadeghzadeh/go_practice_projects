package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"fmt"
)

type Page struct {
	Title string // اسم فایل
	Body  []byte // دیتای فایل
}


var validPath = regexp.MustCompile("^/(edit|save|view|sobhan)/([a-zA-Z0-9]+)$")

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
// ////////////////////////////////////////save//////////////////////////////////////////////////////
func (p *Page) save() error {
	fmt.Println("***save***")

	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}
// //////////////////////////////////////////loadPage////////////////////////////////////////////////////
func loadPage(title string) (*Page, error) {
	fmt.Println("***loadPage***")

	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
// ///////////////////////////////////////////viewHandler///////////////////////////////////////////////////

func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	fmt.Println("***viewHandler***")

	p, err := loadPage(title)

	if err != nil {
		//fmt.Fprintf(w, "<h1>%s</h1>", "File Not Found")
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)

		return
	}
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "view", p)

}
// ////////////////////////////////////////////editHandler//////////////////////////////////////////////////

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("***editHandler***")

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}
// /////////////////////////////////////////////saveHandler/////////////////////////////////////////////////

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("***saveHandler***")

	body := r.FormValue("body")
	//fmt.Println("body: ",body)
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// /////////////////////////////////////////renderTemplate/////////////////////////////////////////////////////

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	fmt.Println("***renderTemplate***")

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Println("***Template Not Found***")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("***Template Not Found-1***")
	}
}

////////////////////////////////////////////getTitle////////////////////////////////////////////////////

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	fmt.Println("***getTitle***")

	m := validPath.FindStringSubmatch(r.URL.Path)
	fmt.Println("m:", m)
	if m == nil {

		//http.NotFound(w, r)
		fmt.Fprint(w, "Invalid Page Title")
		return "", errors.New("Invalid Page Title")
	}
	
	return m[2], nil // The title is the second subexpression.
}
/////////////////////////////////////////////////////////////////////////////////////////////////////
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	fmt.Println("***makeHandler***")
    return func(w http.ResponseWriter, r *http.Request) {
		
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            //http.NotFound(w, r)
			fmt.Fprint(w, "Invalid Route-test")
			fmt.Println("Invalid Route-")
            return
        }
        fn(w, r, m[2])

    }
}
// //////////////////////////////////////////main////////////////////////////////////////////////////
func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

////////////////////////////////////////////////////////////////////////////////////////////////
// func main() {
// 	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
// 	p1.save()
// 	p2, err := loadPage("TestPage")
//     if err!=nil {
//         fmt.Println("File Not Found")
//         return
//     }

// 	fmt.Println(string(p2.Body))
// }
