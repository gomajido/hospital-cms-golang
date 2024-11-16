package dependency

import (
	"github.com/gomajido/hospital-cms-golang/config"
)

type AppUsecase struct {
}

func InitUsecase(config *config.Config, repo *AppRepositories, common *CommonRepositories) *AppUsecase {
	return &AppUsecase{}
}
