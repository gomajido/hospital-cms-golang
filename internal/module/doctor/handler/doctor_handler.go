package handler

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/module/doctor/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
	"github.com/google/uuid"
)

type DoctorHandler struct {
	doctorUsecase domain.DoctorUsecase
}

func NewDoctorHandler(du domain.DoctorUsecase) *DoctorHandler {
	return &DoctorHandler{
		doctorUsecase: du,
	}
}

func (h *DoctorHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateDoctorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	doctor, err := h.doctorUsecase.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(doctor))
}

func (h *DoctorHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	doctor, err := h.doctorUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(doctor))
}

func (h *DoctorHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	doctors, total, err := h.doctorUsecase.List(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	listResponse := response.ListResponse{
		Meta: response.MetaResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Data: doctors,
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(listResponse))
}

func (h *DoctorHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.UpdateDoctorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	doctor, err := h.doctorUsecase.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(doctor))
}

func (h *DoctorHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	if err := h.doctorUsecase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}
	return c.Status(fiber.StatusOK).JSON(response.Ok)

}

func (h *DoctorHandler) CreateSchedule(c *fiber.Ctx) error {
	doctorID, err := uuid.Parse(c.Params("doctor_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	schedule, err := h.doctorUsecase.CreateSchedule(c.Context(), doctorID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(schedule))
}

func (h *DoctorHandler) GetSchedules(c *fiber.Ctx) error {
	doctorID, err := uuid.Parse(c.Params("doctor_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	schedules, err := h.doctorUsecase.GetSchedulesByDoctorID(c.Context(), doctorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(schedules))
}

func (h *DoctorHandler) UpdateSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.UpdateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	schedule, err := h.doctorUsecase.UpdateSchedule(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(schedule))
}

func (h *DoctorHandler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	if err := h.doctorUsecase.DeleteSchedule(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}

func (h *DoctorHandler) CreateReschedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.CreateRescheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	reschedule, err := h.doctorUsecase.CreateReschedule(c.Context(), scheduleID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(reschedule))
}

func (h *DoctorHandler) GetReschedules(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("schedule_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	reschedules, err := h.doctorUsecase.GetReschedulesByScheduleID(c.Context(), scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(reschedules))
}

func (h *DoctorHandler) UpdateReschedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	var req domain.UpdateRescheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	reschedule, err := h.doctorUsecase.UpdateReschedule(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(reschedule))
}

func (h *DoctorHandler) DeleteReschedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(err))
	}

	if err := h.doctorUsecase.DeleteReschedule(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}
