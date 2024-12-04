package dependency

import (
	"github.com/gomajido/hospital-cms-golang/config"
)

type Adapters struct {
}

func InitAdapters(cfg *config.Config, drivers *Drivers) *Adapters {
	return &Adapters{}
}
