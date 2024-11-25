package response

import "net/http"

var (
	Ok = CustomResponse{
		HttpCode: http.StatusOK,
		Message:  "success",
		Code:     "0000",
	}

	SuccessMessage = "success"
)
