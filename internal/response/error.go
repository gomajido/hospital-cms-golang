package response

import (
	"net/http"
)

var (
	ErrInternalServer = CustomResponse{
		HttpCode: http.StatusInternalServerError,
		Message:  "internal server error",
		Code:     "5000",
	}

	ErrBadRequest = CustomResponse{
		HttpCode: http.StatusBadRequest,
		Message:  "bad request",
		Code:     "4000",
	}

	ErrUnauthorized = CustomResponse{
		HttpCode: http.StatusUnauthorized,
		Message:  "Api-Key invalid",
		Code:     "4001",
	}

	ErrInvalidParam = CustomResponse{
		HttpCode: http.StatusBadRequest,
		Message:  "invalid parameter",
		Code:     "4003",
	}

	ErrForbidden = CustomResponse{
		HttpCode: http.StatusForbidden,
		Message:  http.StatusText(http.StatusForbidden),
		Code:     "4003",
	}

	ErrRecordNotFound = CustomResponse{
		HttpCode: http.StatusNotFound,
		Message:  "record not found",
		Code:     "4004",
	}

	ErrUnprocessableEntity = CustomResponse{
		HttpCode: http.StatusUnprocessableEntity,
		Message:  "StatusUnprocessableEntity",
		Code:     "4002",
	}

	StatusTooManyRequests = CustomResponse{
		HttpCode: http.StatusTooManyRequests,
		Message:  "StatusTooManyRequests",
		Code:     "4005",
	}
)
