package api

import (
	"log"
	"net/http"
	"time"
	grpcclients "yaak-kaii/services/api-gateway/internal/grpc_clients"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/auth"

	"github.com/gin-gonic/gin"
)

func HandleCreateUser(ctx *gin.Context) {
	var reqBody types.CreateUserRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		log.Println(err.Error())
		code := http.StatusBadRequest
		ctx.JSON(http.StatusBadRequest, errorResponse(err, code))
		return
	}

	authSvc, err := grpcclients.NewAuthServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer authSvc.Close()

	user, err := authSvc.Client.CreateUser(ctx, &auth.CreateUserRequest{
		Email:       reqBody.Email,
		FirstName:   reqBody.FirstName,
		LastName:    reqBody.LastName,
		PhoneNumber: reqBody.PhoneNumber,
		Password:    reqBody.Password,
	})
	if err != nil {
		code := http.StatusInternalServerError
		ctx.JSON(code, errorResponse(err, code))
		return
	}

	res := contracts.APIResponse{
		Data: types.CreateUserResponse{
			ID:          user.Id,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   user.CreatedAt.AsTime().Format(time.RFC3339),
		},
	}

	ctx.JSON(http.StatusCreated, res)
}
