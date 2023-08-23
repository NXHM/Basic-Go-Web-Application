package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	file_error := openLogFile("Error.txt") // Abre el archivo Error.txt
	defer file_error.Close()               // Cierra el archivo Error.txt al finalizar la ejecución del programa
	log.SetOutput(file_error)              // Escribe los logs en el archivo Error.txt
	http.HandleFunc("/", handler)          // Le dice a http que maneje todas las requests a / con handler
	// Escucha en el puerto 8080 y espera a que termine
	log.Fatal(http.ListenAndServe(":8080", nil)) // Si es que hay un error, se anota el error
}

type Page struct {
	Title string
	Body  []byte
}

func openLogFile(filename string) *os.File {
	// Crea Fila del objeto os.File para escribir logs y err en caso que hayan errores
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) // (nombre del archivo,Agregar contenido|Crear si no existe|Escribir solo, permisos de escritura y lectura para el dueño, grupos y otros usuarios)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
func handler(w http.ResponseWriter, r *http.Request) {
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
