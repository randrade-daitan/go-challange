package grpcApi

import (
	"challange/internal/grpcApi/proto"
	"challange/internal/repository"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type gRPCServer struct {
	db      repository.Repository
	service *grpc.Server
	proto.UnimplementedTaskServiceServer
}

func NewServer(db repository.Repository) *gRPCServer {
	service := grpc.NewServer()
	server := new(gRPCServer)
	server.db = db
	server.service = service
	return server
}

func (server *gRPCServer) StartServing(port int) error {
	err := server.serve(port)
	if err != nil {
		return err
	}

	return server.serveHTTPProxy(port)
}

func (server *gRPCServer) serve(port int) error {
	proto.RegisterTaskServiceServer(server.service, server)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	go func() {
		log.Fatalln(server.service.Serve(listener))
	}()

	return err
}

func (server *gRPCServer) serveHTTPProxy(port int) error {
	mux := runtime.NewServeMux()

	conn, connErr := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%v", port),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if connErr != nil {
		return connErr
	}

	registerErr := proto.RegisterTaskServiceHandler(context.Background(), mux, conn)
	if registerErr != nil {
		return registerErr
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%v", port+1),
		Handler: mux,
	}
	return httpServer.ListenAndServe()
}

func (server *gRPCServer) GetTasks(ctx context.Context, in *proto.Empty) (*proto.Tasks, error) {
	t, err := server.db.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return convertRepoTasks(t), nil
}

func (server *gRPCServer) GetTaskByID(ctx context.Context, in *proto.TaskID) (*proto.Task, error) {
	t, err := server.db.GetTaskByID(in.Id)
	if err != nil {
		return nil, err
	}
	return convertRepoTask(t), nil
}

func (server *gRPCServer) GetTasksByCompletion(ctx context.Context, in *proto.TaskCompletion) (*proto.Tasks, error) {
	c := in.Completed
	if c == nil {
		return server.GetTasks(ctx, &proto.Empty{})
	}

	t, err := server.db.GetTasksByCompletion(*c)
	if err != nil {
		return nil, err
	}
	return convertRepoTasks(t), nil
}

func (server *gRPCServer) AddTask(ctx context.Context, in *proto.Task) (*proto.TaskID, error) {
	newTask := convertDataTask(in)
	id, err := server.db.AddTask(newTask)
	if err != nil {
		return nil, err
	}

	r := &proto.TaskID{
		Id: id,
	}
	return r, nil
}

func (server *gRPCServer) EditTask(ctx context.Context, in *proto.Task) (*proto.Empty, error) {
	newTask := convertDataTask(in)
	err := server.db.EditTask(newTask)
	return &proto.Empty{}, err
}
