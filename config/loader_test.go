package config

import (
	"os"
	"testing"
)

func TestLoader(t *testing.T) {
	loader := NewLoader()

	err := loader.Load()
	if err != nil {
		t.Errorf("can't load .env file: %v", err)
		return
	}

	envs := []string{
		"UID",
		"GID",
		"TRNTLPASS",
		"TRNTLPORT",
		"TRNTLUSER",
	}

	for _, e := range envs {
		if os.Getenv(e) == "" {
			t.Errorf("error occured while loading env. got nothing on %s key", e)
			return
		}
	}

	t.Logf("Passed")
}
