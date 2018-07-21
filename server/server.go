package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/cipepser/gRPC-gacha/gacha"
	"google.golang.org/grpc"
)

type server struct {
}

const (
	port = ":8080"
)

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGachaServer(s, &server{})
	s.Serve(l)
}

func (s *server) Lottery(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	if len(in.Cards) < 1 {
		return nil, errors.New("empty cards")
	}
	rand.Seed(time.Now().UnixNano())
	chosenKey := rand.Intn(len(in.Cards))
	return &pb.Response{Card: in.Cards[chosenKey], RetCode: 1}, nil
}
