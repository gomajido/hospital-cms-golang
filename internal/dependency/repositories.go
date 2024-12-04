package dependency

import (
	"database/sql"

	"github.com/gomajido/hospital-cms-golang/config"
	appointmentDomain "github.com/gomajido/hospital-cms-golang/internal/module/appointment/domain"
	appointmentRepo "github.com/gomajido/hospital-cms-golang/internal/module/appointment/repository"
	articleDomain "github.com/gomajido/hospital-cms-golang/internal/module/article/domain"
	articleRepo "github.com/gomajido/hospital-cms-golang/internal/module/article/repository"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/repository"
	doctorDomain "github.com/gomajido/hospital-cms-golang/internal/module/doctor/domain"
	doctorRepo "github.com/gomajido/hospital-cms-golang/internal/module/doctor/repository"
	"github.com/gomajido/hospital-cms-golang/pkg/db/redis"
)

type CommonRepositories struct {
}

type AppRepositories struct {
	AuthRepo        domain.AuthRepository
	ArticleRepo     articleDomain.ArticleRepository
	DoctorRepo      doctorDomain.DoctorRepository
	AppointmentRepo appointmentDomain.AppointmentRepository
}

func InitCommonRepos(Adapters *Adapters, Drivers *Drivers, config *config.Config) *CommonRepositories {
	return &CommonRepositories{}
}

func InitRepos(db *sql.DB, redis *redis.Redis) *AppRepositories {
	return &AppRepositories{
		AuthRepo:        repository.NewAuthRepository(db),
		ArticleRepo:     articleRepo.NewArticleRepository(db),
		DoctorRepo:      doctorRepo.NewDoctorRepository(db),
		AppointmentRepo: appointmentRepo.NewAppointmentRepository(db),
	}
}
