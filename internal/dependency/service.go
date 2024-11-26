package dependency

import (
	"github.com/gomajido/hospital-cms-golang/config"
	articleDomain "github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
	articleUsecase "github.com/gomajido/hospital-cms-golang/internal/module/article/usecase"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/usecase"
)

type AppUsecase struct {
	AuthUsecase    domain.AuthUsecase
	ArticleUsecase articleDomain.ArticleUsecase
}

func InitUsecase(config *config.Config, repo *AppRepositories, common *CommonRepositories) *AppUsecase {
	return &AppUsecase{
		AuthUsecase:    usecase.NewAuthUsecase(repo.AuthRepo, config),
		ArticleUsecase: articleUsecase.NewArticleUsecase(repo.ArticleRepo),
	}
}
