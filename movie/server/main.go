package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	pb "movie/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Movie struct {
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	ReleaseDate string `json:"releaseDate"`
}

type movieLibraryServer struct {
	pb.UnimplementedMovieLibraryServiceServer // Embed the "unimplemented" gRPC server
	movies                                    []*pb.Movie
}

// u2
func (s *movieLibraryServer) LoadMovies(ctx context.Context, req *pb.MovieRequest) (*pb.MovieResponse, error) {
	// Reset the movie library by overwriting the existing movies.
	s.movies = req.Movies

	JsonFilePath := os.Getenv("JSON_FILE_PATH")
	f, _ := os.Create(JsonFilePath)
	data, _ := json.Marshal(req.Movies)
	_ = os.WriteFile(f.Name(), data, 0644)

	return &pb.MovieResponse{
		StatusCode: http.StatusResetContent,
	}, nil
}

// u3
func (s *movieLibraryServer) GetMovieDetails(ctx context.Context, req *pb.GetMovieDetailsRequest) (*pb.GetMovieDetailsResponse, error) {
	releaseDate := req.ReleaseDate

	JsonFilePath := os.Getenv("JSON_FILE_PATH")
	data, err := os.ReadFile(JsonFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	// Parse the JSON data into a slice of Movie structs
	var movies []Movie
	if err := json.Unmarshal(data, &movies); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Query movies based on the provided release date
	var matchingMovies []*pb.Movie
	for _, movie := range movies {
		if releaseDate == "" || releaseDate == movie.ReleaseDate {
			// Construct a pb.Movie for each matching movie
			matchingMovies = append(matchingMovies, &pb.Movie{
				Title:       movie.Title,
				Genre:       movie.Genre,
				ReleaseDate: movie.ReleaseDate,
			})
		}
	}

	for _, movie := range matchingMovies {
		fmt.Printf("Title: %s, Genre: %s, Release Date: %s\n", movie.Title, movie.Genre, movie.ReleaseDate)
	}

	return &pb.GetMovieDetailsResponse{
		Movies: matchingMovies, // Use the matchingMovies slice
	}, nil
}

// u4
var movieDataMap = map[int32]*pb.Movie{
	1: {
		Title:       "Betty-1",
		Genre:       "crime",
		ReleaseDate: "01-10-2023",
	},
}

func (s *movieLibraryServer) UpdateMovieDetails(ctx context.Context, req *pb.UpdateMovieDetailsRequest) (*pb.UpdateMovieDetailsResponse, error) {
	movieID := req.MovieId

	JsonFilePath := os.Getenv("JSON_FILE_PATH")
	data, err := os.ReadFile(JsonFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	var movies []map[string]interface{}
	if err := json.Unmarshal([]byte(data), &movies); err != nil {
		fmt.Println("Error parsing JSON:", err)

	}

	// Update the second element
	if len(movies) >= 0 {
		movieID = movieID - 1
		movies[movieID]["title"] = req.UpdatedMovie.Title
		movies[movieID]["genre"] = req.UpdatedMovie.Genre
		movies[movieID]["releaseDate"] = req.UpdatedMovie.ReleaseDate
	}

	// Marshal the updated data back to JSON
	updatedJSON, err := json.Marshal(movies)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	f, _ := os.Open(JsonFilePath)
	_ = os.WriteFile(f.Name(), updatedJSON, 0644)

	// Respond with the updated movie
	return &pb.UpdateMovieDetailsResponse{
		StatusCode:   201,
		UpdatedMovie: req.UpdatedMovie,
	}, nil
}

// main
func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	godotenv.Load(".env")

	server := grpc.NewServer()
	pb.RegisterMovieLibraryServiceServer(server, &movieLibraryServer{})

	// Enable reflection for tools like grpcurl
	reflection.Register(server)

	fmt.Println("Movie Library gRPC server started on :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
