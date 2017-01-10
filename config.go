package main

import (
	"github.com/SpirentOrion/luddite"
)

type CloudProvider struct {
	Name     string `yaml:"name"`
	AuthUrl  string `yaml:"auth_url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Tenant   string `yaml:"tenant"`
	Provider string `yaml:"provider"`
}

// Config is a struct that holds config values relevant to the service framework.
type Config struct {
	Service luddite.ServiceConfig
	App     struct {
		Path     string `yaml:"path"`
		Fallback string `yaml:"fallback"`
	}
	CloudProviders struct {
		Providers []CloudProvider `yaml:"providers"`
	}
}
