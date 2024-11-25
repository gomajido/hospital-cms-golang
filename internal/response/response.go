package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	helperValidator "github.com/gomajido/hospital-cms-golang/internal/helper/validator"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
)

type CustomResponse struct {
	HttpCode   int         `json:"http_code"`
	Message    string      `json:"message"`
	Code       string      `json:"code"`
	Detail     string      `json:"detail,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Errors     []ErrorInfo `json:"errors,omitempty"`
}

type ErrorInfo struct {
	Field        string `json:"field"`
	ErrorMessage string `json:"error_message"`
}

func (e CustomResponse) Error() string {
	return e.Message
}

func (e CustomResponse) WithDetail(detail string) CustomResponse {
	customErr := CustomResponse{
		HttpCode: e.HttpCode,
		Message:  e.Message,
		Code:     e.Code,
		Detail:   detail,
	}

	return customErr
}

func (e CustomResponse) WithValidationError(err error) CustomResponse {
	var validationErrors validator.ValidationErrors
	var eMsg string
	errors.As(err, &validationErrors)
	var fields []ErrorInfo
	for _, validationErr := range validationErrors {
		errorString := validationErr.Error()
		validationError, err := helperValidator.ParseValidationError(errorString)
		if err != nil {
			app_log.Errorf("[Response][WithValidationError] Error parse validation error: %v", err)
			eMsg = err.Error()
		} else {
			eMsg = validationError.ErrorMessage
		}
		fields = append(fields, ErrorInfo{
			Field:        validationErr.Field(),
			ErrorMessage: eMsg,
		})
	}
	customErr := CustomResponse{
		HttpCode: e.HttpCode,
		Message:  e.Message,
		Code:     e.Code,
		Detail:   "error validation",
		Errors:   fields,
	}
	return customErr
}

func (e CustomResponse) WithData(data interface{}) CustomResponse {
	customErr := CustomResponse{
		HttpCode: e.HttpCode,
		Message:  e.Message,
		Code:     e.Code,
		Data:     data,
	}

	return customErr
}

func (e CustomResponse) WithPagination(data interface{}, page int, perPage int, total int) CustomResponse {
	if perPage <= 0 {
		perPage = 1
	}

	pagination := map[string]interface{}{
		"current_page": page,
		"per_page":     perPage,
		"total":        total,
		"total_pages":  total / perPage,
	}

	customErr := CustomResponse{
		HttpCode: e.HttpCode,
		Message:  e.Message,
		Code:     e.Code,
		Data:     data,
	}

	customErr.Data = data
	customErr.Pagination = pagination

	return customErr
}

func (e CustomResponse) WithError(err error) CustomResponse {
	errorInfo := ErrorInfo{
		ErrorMessage: err.Error(),
	}

	customErr := CustomResponse{
		HttpCode: e.HttpCode,
		Message:  e.Message,
		Code:     e.Code,
		Detail:   "error validation",
		Errors:   []ErrorInfo{errorInfo},
	}

	return customErr
}

func Build(ctx *fiber.Ctx, ce CustomResponse) {
	respData := map[string]interface{}{
		"code":    ce.Code,
		"message": ce.Message,
	}

	if ce.Data != nil {
		respData["data"] = ce.Data
	}

	if ce.Detail != "" {
		respData["detail"] = ce.Detail
	}

	err := ctx.JSON(respData)

	if err != nil {
		app_log.Errorf("Error build custom response: %v", err)
	}
}

func (e CustomResponse) WithErrorInfo(errInfo []ErrorInfo) CustomResponse {

	customErr := CustomResponse{
		HttpCode: e.HttpCode,
		Message:  e.Message,
		Code:     e.Code,
		Detail:   "error validation",
		Errors:   errInfo,
	}
	return customErr
}
