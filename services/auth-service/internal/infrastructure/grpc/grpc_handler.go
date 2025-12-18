package grpc

import (
	"context"
	"time"
	"yaak-kaii/services/auth-service/internal/domain"
	pb "yaak-kaii/shared/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCHandler struct {
	pb.UnimplementedUserServiceServer
	service domain.UserService
}

func NewGRPCHandler(server *grpc.Server, service domain.UserService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterUserServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	email := req.GetEmail()
	firstName := req.GetFirstName()
	lastName := req.GetLastName()
	phoneNumber := req.GetPhoneNumber()
	password := req.GetPassword()

	arg := &domain.UserModel{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
	}
	err := arg.Password.Set(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the user: %v", err)
	}

	user, err := h.service.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the user: %v", err)
	}

	return &pb.CreateUserResponse{
		Id:            user.ID.String(),
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		PhoneNumber:   user.PhoneNumber,
		IsActive:      user.IsActive,
		RoleId:        user.Role.ID.String(),
		CreatedAt:     timestamppb.New(time.Unix(user.CreatedAt, 0)),
		UpdatedAt:     timestamppb.New(time.Unix(user.UpdatedAt, 0)),
	}, err

}
