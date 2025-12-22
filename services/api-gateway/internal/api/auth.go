package api

import (
	"net/http"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/auth"

	"github.com/gin-gonic/gin"
)

func (app *Application) createUser(ctx *gin.Context) {
	var reqBody types.CreateUserRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		app.responseWithError(ctx, http.StatusBadRequest, err)
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
		app.responseWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	createdAt, _, err := app.formatDate(&user.CreatedAt)
	if err != nil {
		app.responseWithError(ctx, http.StatusInternalServerError, err)
	}

	res := contracts.APIResponse{
		Data: types.CreateUserResponse{
			ID:          user.Id,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   createdAt,
		},
	}

	ctx.JSON(http.StatusCreated, res)
}

func (app *Application) createGuest(ctx *gin.Context) {
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	guest, err := app.GrpcClients.Auth.Client.CreateGuest(ctx, &auth.CreateGuestRequest{
		IpAddress: ipAddress,
		UserAgent: userAgent,
	})

	if err != nil {
		app.responseWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := contracts.APIResponse{
		Data: types.CreateGuestResponse{GuestToken: guest.GuestToken},
	}

	ctx.SetCookie("guest_session", guest.GuestToken, 3600*24*30, "/", "", false, true)
	ctx.JSON(http.StatusCreated, res)
}

func (app *Application) login(ctx *gin.Context) {
	// bind json
	var reqBody types.LoginRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		app.responseWithError(ctx, http.StatusBadRequest, err)
		return
	}

	// login
	loginRes, err := app.GrpcClients.Auth.Client.Login(ctx, &auth.LoginRequest{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})
	if err != nil {
		app.responseWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	// response
	expiresAt, expiresIn, err := app.formatDate(&loginRes.ExpiresAt)
	if err != nil {
		app.responseWithError(ctx, http.StatusInternalServerError, err)
		return
	}
	res := contracts.APIResponse{
		Data: types.LoginResponse{
			UserId:       loginRes.UserId,
			Token:        loginRes.Token,
			RefreshToken: loginRes.RefreshToken,
			ExpiresAt:    expiresAt,
			ExpiresIn:    expiresIn,
		},
	}
	ctx.JSON(http.StatusOK, res)
}
