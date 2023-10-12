package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "movie/proto"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client using the generated package.
	client := pb.NewMovieLibraryServiceClient(conn)

	// Given xml input data need to save in a file through gRPC
	http.HandleFunc("/movie-library/load/", func(w http.ResponseWriter, r *http.Request) {
		// Define the movie records to send.
		movies := []*pb.Movie{
			{
				Title:       "Betty-1",
				Genre:       "crime",
				ReleaseDate: "01-10-2023",
			},
			{
				Title:       "Betty-2",
				Genre:       "sci-fi",
				ReleaseDate: "02-10-2023",
			},
			{
				Title:       "Betty-3",
				Genre:       "drama",
				ReleaseDate: "02-10-2023",
			},
		}

		// Send the movie records to the gRPC service.
		request := &pb.MovieRequest{
			Movies: movies,
		}

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

	})

	// Fetch from xml file
	http.HandleFunc("/movie-library/movie/", func(w http.ResponseWriter, r *http.Request) {
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

		// Serialize the gRPC response to JSON
		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// update
	//u4 code will be implemented

	// Start the HTTP server
	port := 8080
	fmt.Printf("Listening on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
