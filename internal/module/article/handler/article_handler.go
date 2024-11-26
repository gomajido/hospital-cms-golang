package handler

import (
	"errors"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	article, err := h.articleUsecase.GetByID(c.Context(), articleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.JSON(domain.ArticleResponse{
		Data: article,
	})
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
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.JSON(domain.ArticleResponse{
		Data: article,
	})
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
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status")

	if page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("page must be greater than 0")))
	}
	if limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("limit must be greater than 0")))
	}

	articles, total, err := h.articleUsecase.List(c.Context(), page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(domain.ArticlesResponse{
		Data: articles,
		Meta: domain.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  total,
			PerPage:     limit,
		},
	})
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
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uuid.UUID)
	req.AuthorID = userID

	article, err := h.articleUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(domain.ArticleResponse{
		Data: article,
	})
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	var req domain.UpdateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uuid.UUID)

	article, err := h.articleUsecase.Update(c.Context(), articleID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	return c.JSON(domain.ArticleResponse{
		Data: article,
	})
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(errors.New("invalid article ID format")))
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(uuid.UUID)

	err = h.articleUsecase.Delete(c.Context(), articleID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
