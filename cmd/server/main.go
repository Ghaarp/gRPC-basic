package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	generated "github.com/Ghaarp/gRPC-basic/pkg/note_v1"
)

const (
	address = "localhost:"
	port    = "3000"
)

type server struct {
	generated.UnimplementedNoteV1Server
}

type NoteListSynced struct {
	list map[int64]*generated.Note
	mut  sync.RWMutex
}

var noteList NoteListSynced

func (srv *server) Create(context context.Context, request *generated.CreateRequest) (*generated.CreateResponse, error) {

	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()

	note := generated.Note{}

	note.Id = id
	note.Info = request.Info
	note.CreatedAt = timestamppb.New(time.Now())

	noteList.mut.Lock()
	defer noteList.mut.Unlock()

	noteList.list[id] = &note

	response := generated.CreateResponse{}
	response.Id = id
	return &response, nil
}

//func (srv *server) Get(cont context.Context, request *generated.GetRequest) (*generated.GetResponse, error) {

func main() {

	noteList = NoteListSynced{
		list: make(map[int64]*generated.Note),
	}

	listener, err := net.Listen("tcp", address+port)
	if err != nil {
		log.Fatal(err)
	}

	serverObj := grpc.NewServer()
	reflection.Register(serverObj)
	generated.RegisterNoteV1Server(serverObj, &server{})

	log.Printf("Server started om %v", listener.Addr())

	if err := serverObj.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
