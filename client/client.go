package main

import (
	"context"
	"log"

	pb "github.com/cipepser/gRPC-gacha/gacha"
	"google.golang.org/grpc"
)

const (
	address = "localhost"
	port    = ":8080"
)

type client struct {
}

func main() {
	conn, err := grpc.Dial(address+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGachaClient(conn)

	cards := []*pb.Card{
		&pb.Card{Name: "card1"},
		&pb.Card{Name: "card2"},
		&pb.Card{Name: "card3"},
	}

	resp, err := c.Lottery(context.Background(), &pb.Request{Cards: cards})
	if err != nil {
		log.Fatalf("could not get card %v", err)
	}
	log.Printf("get card: %v", resp.Card.Name)
}
