package grpc

import (
	"context"
	"fmt"
	"log"
	"time"
	"yaak-kaii/services/auth-service/internal/domain"
	pb "yaak-kaii/shared/proto/auth"
	"yaak-kaii/shared/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCHandler struct {
	pb.UnimplementedUserServiceServer
	service domain.AuthService
}

func NewGRPCHandler(server *grpc.Server, service domain.AuthService) {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterUserServiceServer(server, handler)
}
func (h *gRPCHandler) VerifyRefreshToken(ctx context.Context, req *pb.VerifyRefreshTokenRequest) (*pb.VerifyRefreshTokenResponse, error) {
	reqRefreshToken := req.GetRefreshToken()

	refreshToken, err := h.service.VerifyRefreshToken(ctx, reqRefreshToken)
	if err != nil {
		return &pb.VerifyRefreshTokenResponse{
			Result: &pb.VerifyRefreshTokenResponse_Error{
				Error: &pb.TokenError{
					Code:    fmt.Sprintln(401),
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &pb.VerifyRefreshTokenResponse{
		Result: &pb.VerifyRefreshTokenResponse_Valid{
			Valid: &pb.TokenValid{UserId: refreshToken.User.ID.String(), ExpiresAt: refreshToken.User.CreatedAt.Unix()},
		},
	}, nil

}
func (h *gRPCHandler) VerifyGuestToken(ctx context.Context, req *pb.VerifyGuestTokenRequest) (*pb.VerifyGuestTokenResponse, error) {
	guestToken := req.GetGuestToken()

	guest, err := h.service.VerifyGuestToken(ctx, guestToken)
	if err != nil {
		log.Println(err)
		return &pb.VerifyGuestTokenResponse{
			Result: &pb.VerifyGuestTokenResponse_Error{
				Error: &pb.TokenError{Code: fmt.Sprint(401), Message: err.Error()},
			},
		}, nil
	}

	return &pb.VerifyGuestTokenResponse{
		Result: &pb.VerifyGuestTokenResponse_Valid{Valid: &pb.TokenValid{UserId: guest.ID.String(), ExpiresAt: guest.ExpiresAt}},
	}, nil

}

func (h *gRPCHandler) CreateGuest(ctx context.Context, req *pb.CreateGuestRequest) (*pb.CreateGuestResponse, error) {
	ipAddress := req.GetIpAddress()
	userAgent := req.GetUserAgent()

	// generate token
	token, tokenHash, err := util.GenerateTokenPair(32)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the guest: %v", err)
	}

	_, err = h.service.CreateGuest(ctx, &domain.GuestModel{
		IpAddress: ipAddress,
		UserAgent: userAgent,
		TokenHash: tokenHash,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the guest: %v", err)
	}

	return &pb.CreateGuestResponse{
		GuestToken: token,
	}, nil
}

func (h *gRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	email := req.GetEmail()
	firstName := req.GetFirstName()
	lastName := req.GetLastName()
	phoneNumber := req.GetPhoneNumber()
	password := req.GetPassword()
	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the user: %v", err)
	}

	arg := &domain.UserModel{
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PhoneNumber:  phoneNumber,
		PasswordHash: passwordHash,
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
