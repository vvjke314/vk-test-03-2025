package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l Loader) Load() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(projectRoot, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		if projectRoot == "/" {
			break
		}

		projectRoot = filepath.Dir(projectRoot)
	}

	return errors.New("no .env file")
}
