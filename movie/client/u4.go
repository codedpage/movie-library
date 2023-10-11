package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "movie/proto"

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

	// Define the movie ID for the movie you want to update.
	movieID := int32(1)

	// Define the updated movie details.
	updatedMovie := &pb.Movie{
		Title:       "Betty-2",
		Genre:       "sci-fi",
		ReleaseDate: "02-10-2023",
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

	// Check the response status code.
	if response.StatusCode == 201 {
		fmt.Println("Movie details updated successfully.")
		fmt.Println("Updated Movie Details:")
		updatedMovieJSON, err := json.Marshal(response.UpdatedMovie)
		if err != nil {
			log.Fatalf("Failed to marshal updated movie details to JSON: %v", err)
		}
		fmt.Println(string(updatedMovieJSON))
	} else {
		log.Printf("Failed to update movie details. HTTP status code: %d\n", response.StatusCode)
	}
}
