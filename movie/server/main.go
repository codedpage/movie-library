package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "movie/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type movieLibraryServer struct {
	pb.UnimplementedMovieLibraryServiceServer // Embed the "unimplemented" gRPC server
	movies                                    []*pb.Movie
}

// u2
func (s *movieLibraryServer) LoadMovies(ctx context.Context, req *pb.MovieRequest) (*pb.MovieResponse, error) {
	// Reset the movie library by overwriting the existing movies.
	s.movies = req.Movies

	fmt.Println(req.Movies)

	return &pb.MovieResponse{
		StatusCode: http.StatusResetContent,
	}, nil
}

// u3
var movieData = map[string]*pb.Movie{
	"01-10-2023": {
		Title:       "Betty-2",
		Genre:       "sci-fi",
		ReleaseDate: "01-10-2023",
	},
	"02-10-2023": {
		Title:       "Betty-3",
		Genre:       "sci-fi",
		ReleaseDate: "02-10-2023",
	},
}

func (s *movieLibraryServer) GetMovieDetails(ctx context.Context, req *pb.GetMovieDetailsRequest) (*pb.GetMovieDetailsResponse, error) {
	releaseDate := req.ReleaseDate

	movie, exists := movieData[releaseDate]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Movie not found for release date: %s", releaseDate)
	}

	return &pb.GetMovieDetailsResponse{
		Movies: []*pb.Movie{movie},
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

	// Check if the movie exists
	if _, exists := movieDataMap[movieID]; !exists {
		return nil, status.Errorf(codes.NotFound, "Movie with ID %d not found", movieID)
	}

	// Update the movie details
	movieDataMap[movieID] = req.UpdatedMovie

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

	server := grpc.NewServer()
	pb.RegisterMovieLibraryServiceServer(server, &movieLibraryServer{})

	// Enable reflection for tools like grpcurl
	reflection.Register(server)

	fmt.Println("Movie Library gRPC server started on :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
