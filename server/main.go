package main

import (
	"context"
	"errors"
	"fmt"
	pb "grpc_tutorial/pb"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	port = ":9000"
)

func main() {

	absoluteFilePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cert := absoluteFilePath + "/cert.pem"
	key := absoluteFilePath + "/key.pem"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error in listening %v\n", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("Error in listening %v\n", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	s := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(s, new(userService))
	log.Println("Service running at port : ", port)
	s.Serve(lis)
}

type userService struct{}

func (s *userService) GetUserById(ctx context.Context, req *pb.UserByIdPayload) (*pb.UserResponse, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Metadata recieved : %v\n", md)
	}

	for _, u := range users {
		if req.Id == u.Id {
			return &pb.UserResponse{User: &u}, nil
		}
	}
	return nil, errors.New("User not found")
}

func (s *userService) GetAllUsers(req *pb.AllUsersPayload, stream pb.UserService_GetAllUsersServer) error {

	for _, u := range users {
		stream.Send(&pb.UserResponse{User: &u})
	}
	return nil
}

func (s *userService) Save(ctx context.Context, req *pb.UserPayload) (*pb.UserResponse, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Metadata recieved : %v\n", md)
	}

	u := new(pb.User)
	u.Id = req.User.Id
	u.Name = req.User.Name
	u.Email = req.User.Email
	u.Password = req.User.Password
	users = append(users, *u)

	return &pb.UserResponse{User: u}, nil
}

func (s *userService) SaveAll(stream pb.UserService_SaveAllServer) error {

	for {
		u, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		users = append(users, *u.User)
		stream.Send(&pb.UserResponse{User: u.User})
	}

	for _, u := range users {
		fmt.Println(u)
	}

	return nil
}
