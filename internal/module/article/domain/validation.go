package domain

import (
	"fmt"

	"github.com/gomajido/hospital-cms-golang/internal/constant"
	"github.com/gomajido/hospital-cms-golang/internal/response"
)

const (
	// Field names for validation messages
	TITLE_FIELD           = "title"
	CONTENT_FIELD         = "content"
	STATUS_FIELD          = "status"
	AUTHOR_ID_FIELD       = "author_id"
	CATEGORY_IDS_FIELD    = "category_ids"
	PAGE_FIELD           = "page"
	LIMIT_FIELD          = "limit"
)

// Validate validates CreateArticleRequest
func (r *CreateArticleRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.Title == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        TITLE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, TITLE_FIELD),
		})
	}

	if r.Content == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        CONTENT_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, CONTENT_FIELD),
		})
	}

	if r.Status == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, STATUS_FIELD),
		})
	} else if r.Status != "published" && r.Status != "draft" && r.Status != "scheduled" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, STATUS_FIELD, "published, draft, or scheduled"),
		})
	}

	if r.AuthorID.String() == "00000000-0000-0000-0000-000000000000" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        AUTHOR_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, AUTHOR_ID_FIELD),
		})
	}

	if len(r.CategoryIDs) == 0 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        CATEGORY_IDS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, CATEGORY_IDS_FIELD),
		})
	} else {
		for i, id := range r.CategoryIDs {
			if id.String() == "00000000-0000-0000-0000-000000000000" {
				errorInfo = append(errorInfo, response.ErrorInfo{
					Field:        fmt.Sprintf("%s[%d]", CATEGORY_IDS_FIELD, i),
					ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, CATEGORY_IDS_FIELD),
				})
			}
		}
	}

	return errorInfo
}

// Validate validates UpdateArticleRequest
func (r *UpdateArticleRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.Status != "" && r.Status != "published" && r.Status != "draft" && r.Status != "scheduled" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, STATUS_FIELD, "published, draft, or scheduled"),
		})
	}

	if len(r.CategoryIDs) > 0 {
		for i, id := range r.CategoryIDs {
			if id.String() == "00000000-0000-0000-0000-000000000000" {
				errorInfo = append(errorInfo, response.ErrorInfo{
					Field:        fmt.Sprintf("%s[%d]", CATEGORY_IDS_FIELD, i),
					ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, CATEGORY_IDS_FIELD),
				})
			}
		}
	}

	return errorInfo
}

// Validate validates ListArticlesRequest
func (r *ListArticlesRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.Page < 1 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        PAGE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_MIN_VALUE, PAGE_FIELD, "1"),
		})
	}

	if r.Limit < 1 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        LIMIT_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_MIN_VALUE, LIMIT_FIELD, "1"),
		})
	} else if r.Limit > 100 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        LIMIT_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_MAX_VALUE, LIMIT_FIELD, "100"),
		})
	}

	if r.Status != "" && r.Status != "published" && r.Status != "draft" && r.Status != "scheduled" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, STATUS_FIELD, "published, draft, or scheduled"),
		})
	}

	return errorInfo
}
