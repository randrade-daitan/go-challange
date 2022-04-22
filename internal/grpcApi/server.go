package grpcApi

import (
	"challange/internal/repository"
	context "context"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type gRPCServer struct {
	db      repository.Repository
	service *grpc.Server
	http.Handler
	UnimplementedTaskServiceServer
}

func NewServer(db repository.Repository) *gRPCServer {
	service := grpc.NewServer()
	server := new(gRPCServer)
	server.db = db
	server.service = service

	RegisterTaskServiceServer(service, server)

	return server
}

func (server *gRPCServer) StartServing(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}
	return server.service.Serve(listener)
}

func (server *gRPCServer) GetAllTasks(context.Context, *emptypb.Empty) (*Tasks, error) {
	t, err := server.db.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return convertRepoTasks(t), nil
}

func (server *gRPCServer) GetTaskByID(_ context.Context, r *TaskID) (*Task, error) {
	t, err := server.db.GetTaskByID(r.Id)
	if err != nil {
		return nil, err
	}
	return convertRepoTask(t), nil
}

func (server *gRPCServer) GetTasksByCompletion(_ context.Context, r *TaskCompletion) (*Tasks, error) {
	t, err := server.db.GetTasksByCompletion(r.Completed)
	if err != nil {
		return nil, err
	}
	return convertRepoTasks(t), nil
}

func (server *gRPCServer) AddTask(_ context.Context, t *Task) (*TaskID, error) {
	newTask := convertDataTask(t)
	id, err := server.db.AddTask(newTask)
	if err != nil {
		return nil, err
	}

	r := &TaskID{
		Id: id,
	}
	return r, nil
}

func (server *gRPCServer) EditTask(_ context.Context, t *Task) (*emptypb.Empty, error) {
	newTask := convertDataTask(t)
	err := server.db.EditTask(newTask)
	return &emptypb.Empty{}, err
}
