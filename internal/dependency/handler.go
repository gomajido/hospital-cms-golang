package dependency

import (
	"context"

	"github.com/gomajido/hospital-cms-golang/config"
)

type ApplicationHandler struct {
}

func InitHandlers(ctx context.Context, cfg *config.Config, service *AppUsecase) *ApplicationHandler {
	return &ApplicationHandler{}
}
