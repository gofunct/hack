package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gofunct/hack/pkg/hackserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

// BarBazServiceServer is a composite interface of foo_pb.BarBazServiceServer and hackserver.Server.
type BarBazServiceServer interface {
	foo_pb.BarBazServiceServer
	hackserver.Server
}

// NewBarBazServiceServer creates a new BarBazServiceServer instance.
func NewBarBazServiceServer() BarBazServiceServer {
	return &barBazServiceServerImpl{}
}

type barBazServiceServerImpl struct {
}

func (s *barBazServiceServerImpl) ListBarBazs(ctx context.Context, req *foo_pb.ListBarBazsRequest) (*foo_pb.ListBarBazsResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *barBazServiceServerImpl) CreateBarBaz(ctx context.Context, req *foo_pb.CreateBarBazRequest) (*foo_pb.BarBaz, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *barBazServiceServerImpl) DeleteBarBaz(ctx context.Context, req *foo_pb.DeleteBarBazRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *barBazServiceServerImpl) Rename(ctx context.Context, req *foo_pb.RenameRequest) (*foo_pb.RenameResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *barBazServiceServerImpl) MoveMove(ctx context.Context, req *foo_pb.MoveMoveRequest) (*foo_pb.MoveMoveResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

