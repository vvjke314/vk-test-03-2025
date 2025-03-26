package config

import "testing"

func TestNewTnConfig(t *testing.T) {
	load := NewLoader()
	load.Load()

	cfg := NewTnConfig()
	if cfg.Username == "" || cfg.Port == "" || cfg.Pass == "" {
		t.Errorf("bad passed arguments")
		return
	}
}
