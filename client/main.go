package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/splax-s/basic-go-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		os.Exit(1) // Exit if the connection failed.
	}
	defer conn.Close()

	// Create a client using the established connection.
	client := pb.NewGreetServiceClient(conn)

	// Define a list of names for the request.
	names := &pb.NamesList{
		Names: []string{"Splax", "Alice", "Bob"},
	}

	// Call the bidirectional streaming function with the client object and list of names.
	callSayHelloBidirectionalStream(client, names)
}

func callSayHelloBidirectionalStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Bidirectional Streaming started")
	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Error while closing stream: %v", err)
	}
	<-waitc
	log.Printf("Bidirectional Streaming finished")
}
