package user_service

import (
	"fmt"
	"log"
	"net"
	"os"

	api "github.com/aburifat/go-agro/apis/agro"
	"github.com/aburifat/go-agro/pkg/backend/services/user_service/handlers"
	"github.com/aburifat/go-agro/pkg/backend/services/user_service/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Server() {
	logger, _ := zap.NewProduction() // or zap.NewDevelopment() for dev
	defer logger.Sync()
	logger.Info("Starting server")
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=users port=5432 sslmode=disable TimeZone=UTC", postgresUser, postgresPassword)
	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	logger.Info("Successfully connected to database")

	//install extension
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	// Auto-migrate the schema (creates/updates tables based on structs)
	err = db.AutoMigrate(&api.User{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	grpcServer := grpc.NewServer()

	proto.RegisterUserServiceServer(grpcServer, handlers.NewUserHandler(db))

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	fmt.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
