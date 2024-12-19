package handlers

import (
	"context"
	"fmt"

	"github.com/aburifat/go-agro/services/user_service/proto/user_proto"
	"github.com/aburifat/go-agro/services/user_service/storage"
)

type UserHandler struct {
	user_proto.UnimplementedUserServiceServer
	storage storage.UserStorage
}

func NewUserHandler(storage storage.UserStorage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *user_proto.CreateUserRequest) (*user_proto.CreateUserResponse, error) {
	user := &storage.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	id, err := h.storage.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &user_proto.CreateUserResponse{
		Id:      id,
		Message: "User created successfully",
	}, nil
}

func (h *UserHandler) GetUserById(ctx context.Context, req *user_proto.GetUserByIdRequest) (*user_proto.GetUserByIdResponse, error) {
	user, err := h.storage.GetById(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user_proto.GetUserByIdResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *UserHandler) GetUsers(ctx context.Context, req *user_proto.GetUsersRequest) (*user_proto.GetUsersResponse, error) {
	users, err := h.storage.GetAll(int(req.GetPageNumber()), int(req.GetPageSize()))
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}

	var userList []*user_proto.User
	for _, u := range users {
		userList = append(userList, &user_proto.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		})
	}

	return &user_proto.GetUsersResponse{
		Users: userList,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *user_proto.UpdateUserRequest) (*user_proto.UpdateUserResponse, error) {
	updatedUser := &storage.User{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	err := h.storage.Update(req.GetId(), updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &user_proto.UpdateUserResponse{
		Message: "User updated successfully",
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *user_proto.DeleteUserRequest) (*user_proto.DeleteUserResponse, error) {
	err := h.storage.Delete(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}

	return &user_proto.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}
