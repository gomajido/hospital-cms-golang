package dependency

import (
	"database/sql"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/pkg/db/redis"
)

type CommonRepositories struct {
}

type AppRepositories struct {
}

func InitCommonRepos(Adapters *Adapters, Drivers *Drivers, config *config.Config) *CommonRepositories {
	return &CommonRepositories{}
}

func InitRepos(db *sql.DB, redis *redis.Redis) *AppRepositories {
	return &AppRepositories{}
}