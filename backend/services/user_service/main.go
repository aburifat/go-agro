package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aburifat/go-agro/services/user_service/handlers"
	"github.com/aburifat/go-agro/services/user_service/proto/user_proto"
	"github.com/aburifat/go-agro/services/user_service/storage"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	mongoStorage, err := storage.NewMongoUserStorage(mongoURI, mongoDBName)
	if err != nil {
		log.Fatalf("failed to initialize MongoDB storage: %v", err)
	}

	grpcServer := grpc.NewServer()

	user_proto.RegisterUserServiceServer(grpcServer, handlers.NewUserHandler(mongoStorage))

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	fmt.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
