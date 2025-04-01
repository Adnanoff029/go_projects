package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type Error struct {
	Message string
}

var movies []Movie

func addMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMovie Movie
	json.NewDecoder(r.Body).Decode(&newMovie)
	newMovie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(newMovie)
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, ele := range movies {
		if ele.ID == params["id"] {
			json.NewEncoder(w).Encode(ele)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
	// w.WriteHeader(http.StatusNotFound)
	// json.NewEncoder(w).Encode(Error{
	// 	Message: "Movie not found",
	// })
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movieData Movie
	json.NewDecoder(r.Body).Decode(&movieData)
	movieData.ID = params["id"]
	for idx, ele := range movies {
		if ele.ID == params["id"] {
			movies[idx] = movieData
		}
	}
	json.NewEncoder(w).Encode(movieData)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var deletedMovie Movie
	for idx, ele := range movies {
		if ele.ID == params["id"] {
			deletedMovie = ele
			movies = append(movies[:idx], movies[idx+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(deletedMovie)
}

func main() {
	movies = append(movies, Movie{ID: "1", Isbn: "435435", Title: "Imaginary", Director: &Director{Firstname: "Faltuu", Lastname: "Mishra"}})
	movies = append(movies, Movie{ID: "2", Isbn: "456445", Title: "Imaginary2", Director: &Director{Firstname: "Faltuu", Lastname: "Mishra"}})
	r := mux.NewRouter()
	r.HandleFunc("/movies/", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}/", getMovie).Methods("GET")
	r.HandleFunc("/movies/", addMovie).Methods("POST")
	r.HandleFunc("/movies/{id}/", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}/", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}

}
