package dependency

import (
	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/usecase"
)

type AppUsecase struct {
	AuthUsecase domain.AuthUsecase
}

func InitUsecase(config *config.Config, repo *AppRepositories, common *CommonRepositories) *AppUsecase {
	return &AppUsecase{
		AuthUsecase: usecase.NewAuthUsecase(repo.AuthRepo, config),
	}
}
