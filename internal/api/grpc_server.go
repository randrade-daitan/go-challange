package api

import (
	"challange/internal/proto"
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

type grpcServer struct {
	repo    repository.Repository
	service *grpc.Server
	proto.UnimplementedTaskServiceServer
}

// Creates a new gRPC/REST combo server.
func NewGrpcServer(repo repository.Repository) *grpcServer {
	service := grpc.NewServer()
	server := new(grpcServer)
	server.repo = repo
	server.service = service
	return server
}

func (server *grpcServer) StartServing(port int) error {
	err := server.serve(port)
	if err != nil {
		return err
	}

	return server.serveHTTPProxy(port)
}

func (server *grpcServer) serve(port int) error {
	proto.RegisterTaskServiceServer(server.service, server)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	go func() {
		log.Fatalln(server.service.Serve(listener))
	}()

	return err
}

func (server *grpcServer) serveHTTPProxy(port int) error {
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

func (server *grpcServer) GetTasks(ctx context.Context, in *proto.Empty) (*proto.Tasks, error) {
	tasks, err := server.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return convertRepoTasks(tasks), nil
}

func (server *grpcServer) GetTaskByID(ctx context.Context, in *proto.TaskID) (*proto.Task, error) {
	task, err := server.repo.GetTaskByID(in.Id)
	if err != nil {
		return nil, err
	}
	return convertRepoTask(task), nil
}

func (server *grpcServer) GetTasksByCompletion(ctx context.Context, in *proto.TaskCompletion) (*proto.Tasks, error) {
	completed := in.Completed
	if completed == nil {
		return server.GetTasks(ctx, &proto.Empty{})
	}

	tasks, err := server.repo.GetTasksByCompletion(*completed)
	if err != nil {
		return nil, err
	}
	return convertRepoTasks(tasks), nil
}

func (server *grpcServer) AddTask(ctx context.Context, in *proto.Task) (*proto.TaskID, error) {
	newTask := convertDataTask(in)
	id, err := server.repo.AddTask(newTask)
	if err != nil {
		return nil, err
	}

	result := &proto.TaskID{
		Id: id,
	}
	return result, nil
}

func (server *grpcServer) EditTask(ctx context.Context, in *proto.Task) (*proto.Empty, error) {
	newTask := convertDataTask(in)
	err := server.repo.EditTask(newTask)
	return &proto.Empty{}, err
}
