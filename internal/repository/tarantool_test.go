package repository

import (
	"context"
	"testing"

	"github.com/vvjke314/vk-test-03-2025/config"
	"github.com/vvjke314/vk-test-03-2025/internal/entities"
)

func initRepository() *TnRepository {
	l := config.NewLoader()
	l.Load()

	ctx := context.Background()
	cfg := config.NewTnConfig()
	repo := NewTnRepository()

	err := repo.Init(ctx, cfg)
	if err != nil {
		return nil
	}

	return repo
}

func TestConnect(t *testing.T) {
	l := config.NewLoader()
	l.Load()

	ctx := context.Background()
	cfg := config.NewTnConfig()
	repo := NewTnRepository()

	err := repo.Init(ctx, cfg)
	if err != nil {
		t.Errorf("can not init repository %v", err)
		return
	}
	defer repo.Close()
}

func TestInsertData(t *testing.T) {
	repo := initRepository()
	defer repo.Close()

	inserts := []entities.VaultItem{
		{"vova", "petya"},
		{"petya", "12"},
		{"13", "vova"},
	}

	var err error
	for i := range inserts {
		err = repo.InsertData(inserts[i])
		if err != nil {
			t.Errorf("error occured while inserting: %v", err)
			return
		}
	}
}

func TestGetAllData(t *testing.T) {
	repo := initRepository()
	defer repo.Close()
	items, err := repo.GetAllData()
	if err != nil {
		t.Fatalf("GetAllData failed: %v", err)
	}

	t.Logf("Got %d ftems", len(items))
	for _, item := range items {
		t.Logf("Key: %s, Value: %s", item.Key, item.Value)
	}
}
