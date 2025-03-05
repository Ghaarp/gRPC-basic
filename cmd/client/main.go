package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	generated "github.com/Ghaarp/gRPC-basic/pkg/note_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
)

const (
	address = "localhost:"
	port    = "3000"
)

func main() {
	connection, err := grpc.Dial(fmt.Sprintf("%s%s", address, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	client := generated.NewNoteV1Client(connection)

	cont, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	noteInfo := generated.NoteInfo{
		Title:    gofakeit.BeerName(),
		Content:  gofakeit.BeerIbu(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	id, err := client.Create(cont, &generated.CreateRequest{Info: &noteInfo})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(color.RedString("Result: %d"), id.Id)

	note, err := client.Get(cont, &generated.GetRequest{Id: id.Id})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(color.RedString("Note: /n", color.GreenString("%v", note)))
}
