package handler

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
	authdomain "github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
)

type ArticleHandler struct {
	articleUsecase domain.ArticleUsecase
}

func NewArticleHandler(au domain.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		articleUsecase: au,
	}
}

// GetByID godoc
// @Summary Get article by ID
// @Description Get article details by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} domain.ArticleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	articleID, err := uuid.Parse(id)
	if err != nil {
		app_log.Errorf("Invalid article ID format: %v, ID: %s", err, id)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	article, err := h.articleUsecase.GetByID(c.Context(), articleID)
	if err != nil {
		app_log.Errorf("Article not found: %v, ID: %s", err, id)
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(article))
}

// GetBySlug godoc
// @Summary Get article by slug
// @Description Get article details by its slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article Slug"
// @Success 200 {object} domain.ArticleResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /articles/slug/{slug} [get]
func (h *ArticleHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	article, err := h.articleUsecase.GetBySlug(c.Context(), slug)
	if err != nil {
		app_log.Errorf("Article not found by slug: %v, slug: %s", err, slug)
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(article))
}

// List godoc
// @Summary List articles
// @Description Get a list of articles with pagination
// @Tags articles
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param status query string false "Article status (published, draft, scheduled)"
// @Success 200 {object} domain.ArticlesResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /articles [get]
func (h *ArticleHandler) List(c *fiber.Ctx) error {
	var req domain.ListArticlesRequest

	// Parse query parameters
	req.Page, _ = strconv.Atoi(c.Query("page", "1"))
	req.Limit, _ = strconv.Atoi(c.Query("limit", "10"))
	req.Status = c.Query("status")

	// Parse category ID if provided
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid category ID format")))
		}
		req.CategoryID = &categoryID
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	articles, total, err := h.articleUsecase.List(c.Context(), req.Page, req.Limit, req.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	listResponse := response.ListResponse{
		Meta: response.MetaResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Data: articles,
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(listResponse))
}

// Create godoc
// @Summary Create a new article
// @Description Create a new article with the provided data
// @Tags articles
// @Accept json
// @Produce json
// @Param article body domain.CreateArticleRequest true "Article data"
// @Success 201 {object} domain.ArticleResponse
// @Failure 400 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /articles [post]
func (h *ArticleHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	// Get user from context
	userToken, ok := c.Locals("user_token").(*authdomain.UserToken)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(fmt.Errorf("missing user token")))
	}

	req.AuthorID = userToken.UserID

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	article, err := h.articleUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Ok.WithData(article))
}

// Update godoc
// @Summary Update an article
// @Description Update an existing article with the provided data
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param article body domain.UpdateArticleRequest true "Article data"
// @Success 200 {object} domain.ArticleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /articles/{id} [put]
func (h *ArticleHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	articleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid article ID format")))
	}

	var req domain.UpdateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	// Validate request
	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	// Get user from context
	userToken, ok := c.Locals("user_token").(*authdomain.UserToken)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(fmt.Errorf("missing user token")))
	}

	// Check if user is author or has permission
	article, err := h.articleUsecase.GetByID(c.Context(), articleID)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	updatedArticle, err := h.articleUsecase.Update(c.Context(), article.ID, userToken.UserID, req)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(updatedArticle))
}

// Delete godoc
// @Summary Delete an article
// @Description Delete an article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /articles/{id} [delete]
func (h *ArticleHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	articleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid article ID format")))
	}

	// Get user from context
	userToken, ok := c.Locals("user_token").(*authdomain.UserToken)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(fmt.Errorf("missing user token")))
	}

	// Check if user is author or has permission
	article, err := h.articleUsecase.GetByID(c.Context(), articleID)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	if err := h.articleUsecase.Delete(c.Context(), article.ID, userToken.UserID); err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusNoContent).JSON(response.Ok)
}

// IncrementVisitorCount godoc
// @Summary Increment article visitor count
// @Description Increment the visitor count of an article
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /articles/{id}/increment-visitor [post]
func (h *ArticleHandler) IncrementVisitorCount(c *fiber.Ctx) error {
	id := c.Params("id")
	articleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid article ID format")))
	}

	if err := h.articleUsecase.IncrementVisitorCount(c.Context(), articleID); err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}
