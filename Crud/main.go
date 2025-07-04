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
	Director *Director `json""director`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = strconv.Itoa(rand.Intn(100000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Movie   Movie  `josn:"movie"`
	}{
		Message: "Movie added successfully",
		Movie:   movie,
	})

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)

	for idx, movie := range movies {

		if movie.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(struct {
		Message string  `json:"message"`
		Movies  []Movie `json:"movies"`
	}{
		Message: "Movie deleted successfully",
		Movies:  movies,
	})

}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	for idx, movie := range movies {

		if movie.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
		}
	}
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Movie   Movie  `json:"movie"`
	}{
		Message: "Movie updated successfully",
		Movie:   movie,
	})

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
