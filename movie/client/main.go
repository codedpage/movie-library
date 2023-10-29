package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	pb "movie/proto"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

// Movie represents a movie with title, genre, and release date.
type Movie struct {
	Title       string `xml:"title,attr"`
	Genre       string `xml:"genre,attr"`
	ReleaseDate string `xml:"releaseDate,attr"`
}

// Movies represents a collection of movies.
type Movies struct {
	XMLName xml.Name `xml:"movies"`
	Movies  []Movie  `xml:"movie"`
}

// MovieLibrary represents a simple movie library.
var MovieLibrary Movies

func resetMovieLibrary() {
	MovieLibrary = Movies{}
}

func loadMovieLibrary(w http.ResponseWriter, r *http.Request) {
	resetMovieLibrary()

	decoder := xml.NewDecoder(r.Body)
	if err := decoder.Decode(&MovieLibrary); err != nil {
		http.Error(w, "Failed to decode XML", http.StatusBadRequest)
		return
	}

	//////////////////////////////////
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client using the generated package.
	client := pb.NewMovieLibraryServiceClient(conn)

	// Send the movie records to the gRPC service.
	movies := make([]*pb.Movie, len(MovieLibrary.Movies))

	for i, v := range MovieLibrary.Movies {
		movies[i] = &pb.Movie{
			Title:       v.Title,
			Genre:       v.Genre,
			ReleaseDate: v.ReleaseDate,
		}
	}

	//fmt.Println(MovieLibrary.Movies)

	request := &pb.MovieRequest{
		Movies: movies,
	}
	fmt.Println(request)

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call the gRPC service's LoadMovies method.
	response, err := client.LoadMovies(ctx, request)
	if err != nil {
		log.Fatalf("Failed to send movie records to gRPC service: %v", err)
	}

	// Check the response status code.
	if response.StatusCode == 205 {
		log.Println("Movie records loaded successfully, and the library is reset.")

	} else {
		log.Printf("Failed to load movie records. HTTP status code: %d\n", response.StatusCode)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusResetContent)
}

func getMovieLibrary(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewMovieLibraryServiceClient(conn)
	releaseDate := r.URL.Path[len("/movie-library/movie/"):]

	// Create a request for querying movie details.
	request := &pb.GetMovieDetailsRequest{
		ReleaseDate: releaseDate,
	}

	// Call the gRPC service's GetMovieDetails method.
	resp, err := client.GetMovieDetails(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to query movie details: %v", err)
	}

	fmt.Println(resp)
	// Serialize the gRPC response to JSON
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func updateMovieLibrary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	decoder := xml.NewDecoder(r.Body)
	if err := decoder.Decode(&MovieLibrary); err != nil {
		http.Error(w, "Failed to decode XML", http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewMovieLibraryServiceClient(conn)

	uriMovieID := r.URL.Path[len("/movie-library/movie/"):]
	parsedID, _ := strconv.ParseInt(uriMovieID, 10, 32)
	movieID := int32(parsedID)
	updatedMovie := &pb.Movie{}

	for _, v := range MovieLibrary.Movies {
		updatedMovie = &pb.Movie{
			Title:       v.Title,
			Genre:       v.Genre,
			ReleaseDate: v.ReleaseDate,
		}
	}

	// Create a request for updating movie details.
	request := &pb.UpdateMovieDetailsRequest{
		MovieId:      movieID,
		UpdatedMovie: updatedMovie,
	}

	// Call the gRPC service's UpdateMovieDetails method.
	response, err := client.UpdateMovieDetails(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to update movie details: %v", err)
	}

	if response.StatusCode == 201 {
		fmt.Println("Movie details updated successfully.")
		fmt.Println("Updated Movie Details:")
		updatedMovieJSON, err := json.Marshal(response.UpdatedMovie)
		if err != nil {
			log.Fatalf("Failed to marshal updated movie details to JSON: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(updatedMovieJSON)
		return

	} else {
		log.Printf("Failed to update movie details. HTTP status code: %d\n", response.StatusCode)
	}
}

func getUpdateMovieLibrary(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		updateMovieLibrary(w, r)
		return
	}

	if r.Method == http.MethodGet {
		getMovieLibrary(w, r)
		return
	}

}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "gRPC api ready!!!")
}


func main() {

	http.HandleFunc("/", apiHandler)
	http.HandleFunc("/movie-library/load", loadMovieLibrary)
	http.HandleFunc("/movie-library/movie/", getUpdateMovieLibrary)
	port := ":8080"
    fmt.Printf("gRPC client is listening on port %s...\n", port)
	http.ListenAndServe(port, nil)
}
