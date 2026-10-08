package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Name     string    `json:"name"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie
var size int = 2

func universalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello!, welcome to my Movie Website")
}
func getALlMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not suppported for this path", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("encoding error: %v", err)
		fmt.Println("Encoding Error\n")
		return
	}

}
func getMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not supported for this path", http.StatusNotFound)
		return
	}
	// query := r.URL.Query()
	// i := query.Get("id")
	// var i int = 0 ;
	vars := mux.Vars(r)
	i := vars["id"]
	var movie Movie
	for _, m := range movies {
		if m.ID == i {
			movie = m
			break
		}
	}
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("encoding error: %v", err)
		fmt.Println("getMovie encoding error\n")
		return
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "method not supported for this path", http.StatusNotFound)
		return
	}
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		return
	}
	for i := range movies {
		if movies[i].ID == movie.ID {
			movies[i] = movie
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}
func createMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported for this path", http.StatusNotFound)
		return
	}
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		return
	}
	movies = append(movies, movie)
	size++
}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not supported for this path", http.StatusNotFound)
		return
	}
	vars := mux.Vars(r)
	i := vars["id"]
	if size == 0 {
		fmt.Println("Not Sufficient size for deleting\n")
		return
	}
	for k := range movies {
		if movies[k].ID == i {
			movies = append(movies[:k], movies[k+1:]...)
			size--
			break
		}
	}
	return
}

// func main() {
// 	fmt.Println("Hello world\n")
// 	router := mux.NewRouter()

// 	movies = append(movies, Movie{ID: "1", Isbn: "438772", Name: "Your Name", Director: &Director{Firstname: "Dunn", Lastname: "smith"}})

// 	movies = append(movies, Movie{ID: "2", Isbn: "438773", Name: "Broken Soul", Director: &Director{Firstname: "Hirata", Lastname: "Ryounuske"}})

// 	router.HandleFunc("/", universalHandler).Methods("GET")
// 	router.HandleFunc("/movies", getALlMovies).Methods("GET")
// 	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
// 	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
// 	router.HandleFunc("/movies", createMovie).Methods("POST")
// 	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

// 	fmt.Println("Server Is Starting On Port 8080\n")
// 	log.Fatal(http.ListenAndServe(":8080", router))

// }
