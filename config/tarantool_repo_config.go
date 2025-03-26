package config

import "os"

type TnRepoConfig struct {
	Username string
	Pass     string
	Port     string
}

func NewTnConfig() *TnRepoConfig {
	return &TnRepoConfig{
		Username: os.Getenv("TRNTLUSER"),
		Pass:     os.Getenv("TRNTLPASS"),
		Port:     os.Getenv("TRNTLPORT"),
	}
}