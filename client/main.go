package main

import (
	"flag"
	"fmt"
	"grpc_tutorial/pb"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
)

const port = ":9000"

func main() {

	absoluteFilePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cert := absoluteFilePath + "/cert.pem"

	option := flag.Int("o", 1, "Command to execute")
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile(cert, "")

	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, opts...)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	switch *option {
	case 1:
		SendMetaData(client)
	case 2:
		GetUserById(client)
	case 3:
		GetAllUsers(client)
	case 4:
		Save(client)
	case 5:
		SaveAll(client)
	default:
		SendMetaData(client)
	}
}

func SendMetaData(client pb.UserServiceClient) {
	md := metadata.MD{}
	md["username"] = []string{"amulya"}
	md["password"] = []string{"itisnotauthenticated"}

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	client.GetUserById(ctx, &pb.UserByIdPayload{})
}

func GetUserById(client pb.UserServiceClient) {
	res, err := client.GetUserById(context.Background(), &pb.UserByIdPayload{Id: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user data : ", res.User)
}

func GetAllUsers(client pb.UserServiceClient) {
	stream, err := client.GetAllUsers(context.Background(), &pb.AllUsersPayload{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("user data is : ", res.User)
	}
}

func Save(client pb.UserServiceClient) {

	new_user := new(pb.User)
	new_user.Id = 3
	new_user.Name = "Amulya3"
	new_user.Email = "amulya3@gmail.com"
	new_user.Password = "imnotwithyouguysanymore"
	res, err := client.Save(context.Background(), &pb.UserPayload{User: new_user})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user data : ", res.User)
}

func SaveAll(client pb.UserServiceClient) {

	users := []pb.User{
		pb.User{Id: 99, Name: "Amulya99", Email: "amulya99@gmail.com", Password: "iamthepassword"},
		pb.User{Id: 98, Name: "Amulya98", Email: "amulya98@gmail.com", Password: "iamthepassword"},
		pb.User{Id: 97, Name: "Amulya97", Email: "amulya97@gmail.com", Password: "iamthepassword"},
	}
	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	doneChannel := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("user data is : ", res.User)
		}
	}()
	for _, u := range users {
		err := stream.Send(&pb.UserPayload{User: &u})
		if err != nil {
			log.Fatal(err)
		}
	}
}
