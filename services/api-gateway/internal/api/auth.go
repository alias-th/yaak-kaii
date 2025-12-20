package api

import (
	"log"
	"net/http"
	"time"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/auth"

	"github.com/gin-gonic/gin"
)

func (app *Application) HandleCreateUser(ctx *gin.Context) {
	var reqBody types.CreateUserRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		log.Println(err.Error())
		code := http.StatusBadRequest
		ctx.JSON(http.StatusBadRequest, errorResponse(err, code))
		return
	}

	user, err := app.GrpcClients.Auth.Client.CreateUser(ctx, &auth.CreateUserRequest{
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

func (app *Application) HandleCreateGuest(ctx *gin.Context) {
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	guest, err := app.GrpcClients.Auth.Client.CreateGuest(ctx, &auth.CreateGuestRequest{
		IpAddress: ipAddress,
		UserAgent: userAgent,
	})

	if err != nil {
		code := http.StatusInternalServerError
		ctx.JSON(code, errorResponse(err, code))
		return
	}

	res := contracts.APIResponse{
		Data: types.CreateGuestResponse{GuestToken: guest.GuestToken},
	}

	ctx.SetCookie("guest_session", guest.GuestToken, 3600*24*30, "/", "", false, true)
	ctx.JSON(http.StatusCreated, res)
}
