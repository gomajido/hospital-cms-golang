package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/module/doctor/domain"
	"github.com/google/uuid"
)

type DoctorHandler struct {
	doctorUsecase domain.DoctorUsecase
}

func NewDoctorHandler(app *fiber.App, du domain.DoctorUsecase) {
	handler := &DoctorHandler{
		doctorUsecase: du,
	}

	api := app.Group("/api/v1")
	doctors := api.Group("/doctors")

	doctors.Post("/", handler.Create)
	doctors.Get("/:id", handler.GetByID)
	doctors.Put("/:id", handler.Update)
	doctors.Delete("/:id", handler.Delete)
	doctors.Get("/", handler.List)

	// Schedule routes
	doctors.Post("/:id/schedules", handler.CreateSchedule)
	doctors.Get("/:id/schedules", handler.GetSchedules)
	doctors.Put("/schedules/:id", handler.UpdateSchedule)
	doctors.Delete("/schedules/:id", handler.DeleteSchedule)

	// Reschedule routes
	doctors.Post("/schedules/:id/reschedules", handler.CreateReschedule)
	doctors.Get("/schedules/:id/reschedules", handler.GetReschedules)
	doctors.Put("/reschedules/:id", handler.UpdateReschedule)
	doctors.Delete("/reschedules/:id", handler.DeleteReschedule)
}

func (h *DoctorHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateDoctorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	doctor, err := h.doctorUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(doctor)
}

func (h *DoctorHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	doctor, err := h.doctorUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Doctor not found",
		})
	}

	return c.JSON(doctor)
}

func (h *DoctorHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	doctors, total, err := h.doctorUsecase.List(c.Context(), page, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": doctors,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *DoctorHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req domain.UpdateDoctorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	doctor, err := h.doctorUsecase.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(doctor)
}

func (h *DoctorHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	if err := h.doctorUsecase.Delete(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Schedule handlers
func (h *DoctorHandler) CreateSchedule(c *fiber.Ctx) error {
	doctorID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid doctor ID format",
		})
	}

	var req domain.CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	schedule, err := h.doctorUsecase.CreateSchedule(c.Context(), doctorID, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(schedule)
}

func (h *DoctorHandler) GetSchedules(c *fiber.Ctx) error {
	doctorID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid doctor ID format",
		})
	}

	schedules, err := h.doctorUsecase.GetSchedulesByDoctorID(c.Context(), doctorID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(schedules)
}

func (h *DoctorHandler) UpdateSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID format",
		})
	}

	var req domain.UpdateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	schedule, err := h.doctorUsecase.UpdateSchedule(c.Context(), id, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(schedule)
}

func (h *DoctorHandler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID format",
		})
	}

	if err := h.doctorUsecase.DeleteSchedule(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Reschedule handlers
func (h *DoctorHandler) CreateReschedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID format",
		})
	}

	var req domain.CreateRescheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	reschedule, err := h.doctorUsecase.CreateReschedule(c.Context(), scheduleID, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(reschedule)
}

func (h *DoctorHandler) GetReschedules(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schedule ID format",
		})
	}

	reschedules, err := h.doctorUsecase.GetReschedulesByScheduleID(c.Context(), scheduleID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(reschedules)
}

func (h *DoctorHandler) UpdateReschedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reschedule ID format",
		})
	}

	var req domain.UpdateRescheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	reschedule, err := h.doctorUsecase.UpdateReschedule(c.Context(), id, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(reschedule)
}

func (h *DoctorHandler) DeleteReschedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reschedule ID format",
		})
	}

	if err := h.doctorUsecase.DeleteReschedule(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
