package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/module/appointment/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
	"github.com/google/uuid"
)

type AppointmentHandler struct {
	appointmentUsecase domain.AppointmentUsecase
}

func NewAppointmentHandler(au domain.AppointmentUsecase) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentUsecase: au,
	}
}

func (h *AppointmentHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	// Get user ID from authenticated context
	userID := c.Locals("user_id").(uuid.UUID)
	req.UserID = userID

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	appointment, err := h.appointmentUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(appointment))
}

func (h *AppointmentHandler) Cancel(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.CancelAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	// Get user ID from authenticated context
	userID := c.Locals("user_id").(uuid.UUID)
	req.UserID = userID

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	appointment, err := h.appointmentUsecase.Cancel(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(appointment))
}

func (h *AppointmentHandler) Reschedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.RescheduleAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	// Get user ID from authenticated context
	userID := c.Locals("user_id").(uuid.UUID)
	req.UserID = userID

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	appointment, err := h.appointmentUsecase.Reschedule(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(appointment))
}

func (h *AppointmentHandler) CheckAvailability(c *fiber.Ctx) error {
	var req domain.CheckAvailabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	isAvailable, err := h.appointmentUsecase.CheckAvailability(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(fiber.Map{
		"is_available": isAvailable,
	}))
}

func (h *AppointmentHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	appointment, err := h.appointmentUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(appointment))
}

func (h *AppointmentHandler) GetByUserID(c *fiber.Ctx) error {
	// Get user ID from authenticated context
	userID := c.Locals("user_id").(uuid.UUID)

	appointments, totalCount, err := h.appointmentUsecase.GetByUserID(c.Context(), userID, 1, 10) // Adding default pagination
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(fiber.Map{
		"appointments": appointments,
		"total_count":  totalCount,
	}))
}

func (h *AppointmentHandler) GetByDoctorID(c *fiber.Ctx) error {
	doctorID, err := uuid.Parse(c.Params("doctor_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	appointments, totalCount, err := h.appointmentUsecase.GetByDoctorID(c.Context(), doctorID, 1, 10) // Adding default pagination
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(fiber.Map{
		"appointments": appointments,
		"total_count":  totalCount,
	}))
}
