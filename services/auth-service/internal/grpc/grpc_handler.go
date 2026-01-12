package grpc

import (
	"context"
	"log"
	"yaak-kaii/services/auth-service/internal/domain"
	pb "yaak-kaii/shared/proto/auth"
	"yaak-kaii/shared/utils"

	"google.golang.org/grpc"
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

func (h *gRPCHandler) VerifyRefreshToken(ctx context.Context, req *pb.VerifyRefreshTokenRequest) (*pb.VerifyTokenResponse, error) {
	reqRefreshToken := req.GetRefreshToken()

	refreshToken, err := h.service.VerifyRefreshToken(ctx, reqRefreshToken)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Result: &pb.VerifyTokenResponse_Error{
				Error: &pb.ErrorResponse{
					Code:    pb.ErrorCode_TOKEN_EXPIRED,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Result: &pb.VerifyTokenResponse_Valid{
			Valid: &pb.TokenValid{
				UserId:    refreshToken.User.ID.String(),
				ExpiresAt: refreshToken.ExpiresAt,
			},
		},
	}, nil

}
func (h *gRPCHandler) VerifyGuestToken(ctx context.Context, req *pb.VerifyGuestTokenRequest) (*pb.VerifyTokenResponse, error) {
	guestToken := req.GetGuestToken()

	guest, err := h.service.VerifyGuestToken(ctx, guestToken)
	if err != nil {
		log.Println(err)
		return &pb.VerifyTokenResponse{
			Result: &pb.VerifyTokenResponse_Error{
				Error: &pb.ErrorResponse{
					Code:    pb.ErrorCode_TOKEN_EXPIRED,
					Message: err.Error()},
			},
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Result: &pb.VerifyTokenResponse_Valid{Valid: &pb.TokenValid{UserId: guest.ID.String(), ExpiresAt: guest.ExpiresAt}},
	}, nil

}

func (h *gRPCHandler) CreateGuest(ctx context.Context, req *pb.CreateGuestRequest) (*pb.CreateGuestResponse, error) {
	ipAddress := req.GetIpAddress()
	userAgent := req.GetUserAgent()

	token, err := h.service.CreateGuest(ctx, &domain.CreateGuestRequest{
		IpAddress: ipAddress,
		UserAgent: userAgent,
	})

	if err != nil {
		return nil, utils.ToGRPCError(err)
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

	// 1. Map data
	arg := &domain.CreateUserRequest{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	// 2. Create user
	user, err := h.service.CreateUser(ctx, arg)

	// 3. Check errors
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	// 4. response
	return &pb.CreateUserResponse{
		Id:            user.User.ID.String(),
		Email:         user.User.Email,
		EmailVerified: user.User.EmailVerified,
		FirstName:     user.User.FirstName,
		LastName:      user.User.LastName,
		PhoneNumber:   user.User.PhoneNumber,
		IsActive:      user.User.IsActive,
		RoleId:        user.User.Role.ID.String(),
		CreatedAt:     user.User.CreatedAt,
		UpdatedAt:     user.User.UpdatedAt,
	}, nil

}

func (h *gRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	arg := &domain.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := h.service.Login(ctx, arg)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.LoginResponse{
		UserId:       resp.UserID,
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	}, nil
}

func (h *gRPCHandler) RotateRefreshToken(
	ctx context.Context, req *pb.RotateRefreshTokenRequest,
) (*pb.RotateRefreshTokenResponse, error) {
	refreshToken := req.GetRefreshToken()

	resp, err := h.service.RotateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, utils.ToGRPCError(err)
	}

	return &pb.RotateRefreshTokenResponse{
		UserId:       resp.UserID,
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt,
	}, nil
}
