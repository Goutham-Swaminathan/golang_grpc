package main

import (
	pb "grpc_tutorial/pb"
)

//this is temp database 
//it should be replaced with any DB data
var users = []pb.User{
	pb.User{
		Id:       1,
		Name:     "Amulya1",
		Email:    "amulya1@gmail.com",
		Password: "imdonenow1",
	},
	pb.User{
		Id:       2,
		Name:     "Amulya2",
		Email:    "amulya2@gmail.com",
		Password: "imdonenow2",
	},
}
