package handler

import (
	"errors"
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
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		app_log.Errorf("Invalid page number: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		app_log.Errorf("Invalid limit number: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	status := c.Query("status")

	if page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("page must be greater than 0")))
	}
	if limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("limit must be greater than 0")))
	}

	articles, total, err := h.articleUsecase.List(c.Context(), page, limit, status)
	if err != nil {
		app_log.Errorf("Failed to fetch articles: %v, page: %d, limit: %d", err, page, limit)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(domain.ArticlesResponse{
		Data: articles,
		Meta: domain.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  total,
			PerPage:     limit,
		},
	}))
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
		app_log.Errorf("Failed to parse create article request: %v", err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Get user ID from user token in context
	userToken, ok := c.Locals("user_token").(*authdomain.UserToken)
	if !ok || userToken == nil {
		app_log.Error("User token not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrUnauthorized)
	}

	userID, err := uuid.Parse(userToken.UserID)
	if err != nil {
		app_log.Errorf("Invalid user ID format in token: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid user ID format")))
	}

	req.AuthorID = userID

	article, err := h.articleUsecase.Create(c.Context(), req)
	if err != nil {
		app_log.Errorf("Failed to create article: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}
	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(article))
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
		app_log.Errorf("Invalid article ID format: %v, ID: %s", err, id)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	var req domain.UpdateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		app_log.Errorf("Failed to parse update article request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	// Get user ID from user token in context
	userToken, ok := c.Locals("user_token").(*authdomain.UserToken)
	if !ok || userToken == nil {
		app_log.Error("User token not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrUnauthorized)
	}

	userID, err := uuid.Parse(userToken.UserID)
	if err != nil {
		app_log.Errorf("Invalid user ID format in token: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid user ID format")))
	}

	article, err := h.articleUsecase.Update(c.Context(), articleID, userID, req)
	if err != nil {
		app_log.Errorf("Failed to update article: %v, ID: %s, userID: %s", err, id, userID)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(article))
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
		app_log.Errorf("Invalid article ID format: %v, ID: %s", err, id)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	// Get user ID from user token in context
	userToken, ok := c.Locals("user_token").(*authdomain.UserToken)
	if !ok || userToken == nil {
		app_log.Error("User token not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrUnauthorized)
	}

	userID, err := uuid.Parse(userToken.UserID)
	if err != nil {
		app_log.Errorf("Invalid user ID format in token: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid user ID format")))
	}

	err = h.articleUsecase.Delete(c.Context(), articleID, userID)
	if err != nil {
		app_log.Errorf("Failed to delete article: %v, ID: %s, userID: %s", err, id, userID)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
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
		app_log.Errorf("Invalid article ID format: %v, ID: %s", err, id)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	err = h.articleUsecase.IncrementVisitorCount(c.Context(), articleID)
	if err != nil {
		app_log.Errorf("Failed to increment visitor count: %v, ID: %s", err, id)
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}
