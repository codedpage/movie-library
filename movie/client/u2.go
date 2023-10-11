package main

import (
	"context"
	"log"
	"time"

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
}
