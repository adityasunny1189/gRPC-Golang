package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "moviesapp.com/grpc/protos"
)

const (
	addr = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()
	client := pb.NewMovieClient(conn)

	runGetMovies(client)
	// runGetMovie(client, "1")
	// runCreateMovie(client, "0736484", "Dangal", "Amir", "Khan")
	// runUpdateMovie(client, "1", "278384", "Bahubali", "S.S", "Rajamouli")
	// runDeleteMovie(client, "98081")

}

func runGetMovies(client pb.MovieClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetMovies(ctx, req)
	if err != nil {
		log.Fatalf("error fetching movies: %v", err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error %v %v", client, err)
		}
		fmt.Println("MoviesInfo: ", row)
	}
}

func runGetMovie(client pb.MovieClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.ID{Value: movieid}
	res, err := client.GetMovie(ctx, req)
	if err != nil {
		log.Printf("error %v %v", client, err)
	}
	fmt.Println("MoviesInfo: ", res)
}

func runCreateMovie(client pb.MovieClient, isbn, title, firstname, lastname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.MoviesInfo{
		Isbn:  isbn,
		Title: title,
		Director: &pb.Director{
			Firstname: firstname,
			Lastname:  lastname,
		},
	}
	res, err := client.CreateMovie(ctx, req)
	if err != nil {
		log.Printf("error %v %v", client, err)
	}
	if res.GetValue() != "" {
		log.Printf("Created Movie Id: %v", res)
	} else {
		log.Printf("Create movie failed")
	}
}

func runUpdateMovie(client pb.MovieClient, movieid, isbn, title, firstname, lastname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.MoviesInfo{
		Id:    movieid,
		Isbn:  isbn,
		Title: title,
		Director: &pb.Director{
			Firstname: firstname,
			Lastname:  lastname,
		},
	}
	res, err := client.UpdateMovie(ctx, req)
	if err != nil {
		log.Printf("error %v %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("Updated movie success: %v", res)
	} else {
		log.Printf("update movie failed")
	}
}

func runDeleteMovie(client pb.MovieClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.ID{Value: movieid}
	res, err := client.DeleteMovie(ctx, req)
	if err != nil {
		log.Printf("error %v %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("delete movie successful: %v", res)
	} else {
		log.Printf("delete movie failed")
	}
}
