package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func main() {
	file_error := writeErrorsFile("Error.txt") // Abre el archivo Error.txt
	defer file_error.Close()                   // Cierra el archivo Error.txt al finalizar la ejecución del programa
	log.SetOutput(file_error)                  // Escribe los logs en el archivo Error.txt
	// http.HandleFunc("/", handler)          // Le dice a http que maneje todas las requests a / con handler
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	// Uso select case segun el path que se coloque en el navegador
	err := http.ListenAndServe(":8080", nil)
	// En caso que haya un error, se anota el error
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		os.Exit(1)
	}
}

type Page struct {
	Title string
	Body  []byte
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, p)
}
func writeErrorsFile(filename string) *os.File {
	// Crea Fila del objeto os.File para escribir logs y err en caso que hayan errores
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) // (nombre del archivo,Agregar contenido|Crear si no existe|Escribir solo, permisos de escritura y lectura para el dueño, grupos y otros usuarios)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
func handler(w http.ResponseWriter, r *http.Request) { // w es la respuesta que se le envía al cliente, r es el request que se recibe del cliente y r es la solicitud
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
