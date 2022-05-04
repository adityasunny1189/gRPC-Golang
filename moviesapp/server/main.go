package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	"google.golang.org/grpc"
	pb "moviesapp.com/grpc/protos"
)

const (
	port = ":8080"
)

var movies []*pb.MoviesInfo

type movieServer struct {
	pb.UnimplementedMovieServer
}

func main() {
	initMovies()

	lst, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMovieServer(s, &movieServer{})
	if err := s.Serve(lst); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initMovies() {
	movie1 := &pb.MoviesInfo{
		Id:    "1",
		Isbn:  "0567834",
		Title: "KGF",
		Director: &pb.Director{
			Firstname: "Neel",
			Lastname:  "Kumar",
		},
	}
	movies = append(movies, movie1)
}

func (s *movieServer) GetMovies(in *pb.Empty, stream pb.Movie_GetMoviesServer) error {
	log.Printf("Received: %v", in)
	for _, movie := range movies {
		if err := stream.Send(movie); err != nil {
			return err
		}
	}
	return nil
}

func (s *movieServer) GetMovie(ctx context.Context, in *pb.ID) (*pb.MoviesInfo, error) {
	log.Printf("Received: %v", in)
	res := &pb.MoviesInfo{}

	for _, movie := range movies {
		if movie.GetId() == in.GetValue() {
			res = movie
			break
		}
	}
	return res, nil
}

func (s *movieServer) CreateMovie(ctx context.Context, in *pb.MoviesInfo) (*pb.ID, error) {
	log.Printf("Received: %v", in)
	res := pb.ID{}
	res.Value = strconv.Itoa(rand.Intn(100000))
	in.Id = res.GetValue()
	movies = append(movies, in)
	return &res, nil
}

func (s *movieServer) UpdateMovie(ctx context.Context, in *pb.MoviesInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	for i, movie := range movies {
		if movie.GetId() == in.GetId() {
			movies = append(movies[:i], movies[i+1:]...)
			in.Id = movie.Id
			movies = append(movies, in)
			res.Value = 1
			break
		}
	}
	return &res, nil
}

func (s *movieServer) DeleteMovie(ctx context.Context, in *pb.ID) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	for i, movie := range movies {
		if movie.GetId() == in.GetValue() {
			movies = append(movies[:i], movies[i+1:]...)
			res.Value = 1
			break
		}
	}
	return &res, nil
}
