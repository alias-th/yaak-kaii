package api

import (
	"errors"
	"log"
	"net/http"
	"time"
	"yaak-kaii/shared/contracts"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (app *Application) responseWithError(ctx *gin.Context, code int, err error) {
	log.Printf("Error: %v", err)

	// function จาก gRPC package ที่ใช้แยกข้อมูล gRPC status จาก error
	if st, ok := status.FromError(err); ok {
		httpCode := mapGrpcCodeToHTTP(st.Code())
		ctx.JSON(httpCode, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    st.Code().String(),
				Message: st.Message(),
			},
		})
		return
	}

	// Fallback
	ctx.JSON(code, contracts.APIResponse{
		Error: &contracts.APIError{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	})
}

func (app *Application) responseWithValidationError(ctx *gin.Context, err error) bool {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return false
	}

	log.Printf("Validation errors: %v", validationErrs)

	ctx.JSON(http.StatusBadRequest, contracts.APIResponse{
		Error: &contracts.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "invalid_request",
		},
	})
	return true
}

func mapGrpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict // 409
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (app *Application) formatDate(timestamp *int64) (at string, in int64, err error) {
	if timestamp == nil {
		return "", 0, errors.New("expires_at timestamp is required")
	}

	expiresAt := time.Unix(*timestamp, 0)
	expiresIn := int64(time.Until(expiresAt).Seconds())

	return expiresAt.Format(time.RFC3339), expiresIn, nil
}

func isAllowedImageCT(ct string) bool {
	switch ct {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func guessExtFromContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
