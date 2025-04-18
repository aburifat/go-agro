package handlers

import (
	"context"
	"fmt"

	api "github.com/aburifat/go-agro/apis/agro"
	"github.com/aburifat/go-agro/pkg/backend/services/user_service/proto"
	"github.com/aburifat/go-agro/pkg/backend/services/user_service/repository"

	"gorm.io/gorm"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	userHandler := UserHandler{
		db: db,
	}
	return &userHandler
}

func (h *UserHandler) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user := &api.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	id, err := repository.Create(h.db, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &proto.CreateUserResponse{
		Id:      id,
		Message: "User created successfully",
	}, nil
}

func (h *UserHandler) GetUserById(ctx context.Context, req *proto.GetUserByIdRequest) (*proto.GetUserByIdResponse, error) {
	user, err := repository.GetById(h.db, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &proto.GetUserByIdResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *UserHandler) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	users, err := repository.GetAll[api.User](h.db, int(req.GetPageNumber()), int(req.GetPageSize()))
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	var userList []*proto.User
	for _, u := range users {
		userList = append(userList, &proto.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		})
	}

	return &proto.GetUsersResponse{
		Users: userList,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	updatedUser := &api.User{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	err := repository.Update(h.db, req.GetId(), updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &proto.UpdateUserResponse{
		Message: "User updated successfully",
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := repository.Delete(h.db, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}
	return &proto.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}
