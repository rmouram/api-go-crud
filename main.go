package main

/*
video base: https://www.youtube.com/watch?v=7W50SBMs6iI&list=PLUbb2i4BuuzDFzmdGK00pLddrIBO6JWyV&index=5
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Book ...
type Book struct {
	Id     int    `json:"id"` //para q o json esteja em minusculo
	Author string `json:"author"`
	Title  string `json:"title"`
}

// Books ...
var Books []Book = []Book{
	{
		Id:     1,
		Author: "Jose de Alencar",
		Title:  "O guarani",
	}, {
		Id:     2,
		Author: "Jose de Alencar",
		Title:  "Iracema",
	}, {
		Id:     3,
		Author: "Jose Saramago",
		Title:  "Ensaio sobre a Cegueira",
	},
}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	urlSplit := strings.Split(r.URL.Path, "/")

	var indiceBook = -1
	id, _ := strconv.Atoi(urlSplit[2])
	for indice, book := range Books {
		if book.Id == id {
			indiceBook = indice
			break
		}
	}
	if indiceBook < 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	Books = append(Books[0:indiceBook], Books[indiceBook+1:len(Books)]...)

	w.WriteHeader(http.StatusNoContent)
}

func searchBooks(w http.ResponseWriter, r *http.Request) {

	enableCors(w)

	w.Header().Set("Content-Type", "application/json")
	//e.g. /livros/123 ---> ["","livros","123"]
	urlSplit := strings.Split(r.URL.Path, "/")

	if len(urlSplit) >= 4 && urlSplit[3] != "" {
		fmt.Println("com barra pode")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//var id, _ = strconv.Atoi(urlSplit[2])
	var title = strings.ToLower(urlSplit[2])
	var encontrou = 0
	for _, livro := range Books {
		if strings.ToLower(livro.Title) == title {
			encontrou = 1
			json.NewEncoder(w).Encode(livro)
		}
	}
	if encontrou == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	w.Header().Set("Content-Type", "application/json")
	//transforma a resposta em um formato json
	json.NewEncoder(w).Encode(Books)
}

func registerBooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enableCors(w)
	indexHandler(w, r)
	// ioutil.ReadAll ler o que é passado como parametro
	// r.Body acessa o conteúdo do body da req
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var newBook Book
	json.Unmarshal(body, &newBook)

	// criar um id automaticamente
	newBook.Id = len(Books) + 1

	Books = append(Books, newBook)

	json.NewEncoder(w).Encode(newBook)
}

func routeBooks(w http.ResponseWriter, r *http.Request) {
	indexHandler(w, r)
	enableCors(w)
	// /books/123/
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) == 2 || len(parts) == 3 && parts[2] == "" {
		if r.Method == "GET" {
			listBooks(w, r)
		} else if r.Method == "POST" {
			registerBooks(w, r)
		}
	} else if len(parts) == 3 || len(parts) == 4 && parts[3] == "" {
		if r.Method == "GET" {
			searchBooks(w, r)
		} else if r.Method == "DELETE" {
			deleteBooks(w, r)
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

}

// as proximas 3 funções são para habilitar o CORS
func enableCors(w http.ResponseWriter) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupResponse(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	setupResponse(w, req)
	if (*req).Method == "OPTIONS" {
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	http.ServeFile(w, r, "index.html")
}

func handlerConfig() {
	http.HandleFunc("/", home)
	http.HandleFunc("/books", routeBooks)
	http.HandleFunc("/books/", routeBooks)
}

func serverConfig() {
	fmt.Println("Servidor esta rodando!")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handlerConfig()
	serverConfig()
}
