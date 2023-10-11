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

	// Define the release date for the movie(s) you want to query.
	releaseDate := "02-10-2023"

	// Create a request for querying movie details.
	request := &pb.GetMovieDetailsRequest{
		ReleaseDate: releaseDate,
	}

	// Call the gRPC service's GetMovieDetails method.
	response, err := client.GetMovieDetails(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to query movie details: %v", err)
	}

	// Convert the response to JSON.
	movieDetailsJSON, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Failed to marshal movie details to JSON: %v", err)
	}

	fmt.Println("Movie Details:")
	fmt.Println(string(movieDetailsJSON))
}
