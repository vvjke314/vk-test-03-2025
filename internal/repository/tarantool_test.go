package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/vvjke314/vk-test-03-2025/config"
	"github.com/vvjke314/vk-test-03-2025/internal/entities"
)

type MockLogger struct {
}

func (m MockLogger) Error(message string) {
	fmt.Println(message)
}

func (m MockLogger) Info(message string) {
	fmt.Println(message)
}

func initRepository() *TnRepository {
	l := config.NewLoader()
	l.Load()

	ctx := context.Background()
	cfg := config.NewTnConfig()
	repo := NewTnRepository()

	err := repo.Init(ctx, cfg, MockLogger{})
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

	err := repo.Init(ctx, cfg, MockLogger{})
	if err != nil {
		t.Errorf("can not init repository %v", err)
		return
	}
	defer repo.Close()
}

func TestInsert(t *testing.T) {
	repo := initRepository()
	defer repo.Close()

	inserts := []entities.VaultItem{
		{"vova", "petya"},
		{"petya", "12"},
		{"13", "vova"},
	}

	var err error
	for i := range inserts {
		err = repo.Insert(inserts[i])
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

	t.Logf("Got %d items", len(items))
	for _, item := range items {
		t.Logf("Key: %s, Value: %s", item.Key, item.Value)
	}
}

type KeyExistenceCase struct {
	key   string
	exist bool
}

func TestKeyExistence(t *testing.T) {
	repo := initRepository()
	defer repo.Close()

	var keys = []KeyExistenceCase{
		{"13", true},
		{"hello", false},
		{"petya", true},
		{"vova", true},
		{"Vova", false},
	}

	for _, k := range keys {
		exist, err := repo.KeyExists(k.key)
		if err != nil {
			t.Errorf("error occured. %v", err)
			return
		}
		if k.exist != exist {
			t.Errorf("error occured. wait %v got %v", k.exist, exist)
			return
		}
	}
}

func TestDelete(t *testing.T) {
	repo := initRepository()
	defer repo.Close()

	data := entities.VaultItem{"hello", "world"}

	err := repo.Insert(data)
	if err != nil {
		t.Errorf("failed while inserting data: %v", err)
		return
	}

	err = repo.Delete(data.Key)
	if err != nil {
		t.Errorf("failed while deleting data: %v", err)
		return
	}
}

func TestUpdate(t *testing.T) {
	repo := initRepository()
	defer repo.Close()

	data := entities.VaultItem{"hello", "world"}

	err := repo.Insert(data)
	if err != nil {
		t.Errorf("failed while inserting data: %v", err)
		return
	}
	defer repo.Delete(data.Key)

	err = repo.Update(data.Key, "tarantool!")
	if err != nil {
		t.Errorf("failed while deleting data: %v", err)
		return
	}
}

func TestGet(t *testing.T) {
	repo := initRepository()
	defer repo.Close()

	data := entities.VaultItem{
		Key:   "hello",
		Value: "world",
	}

	err := repo.Insert(data)
	if err != nil {
		t.Errorf("failed while inserting data: %v", err)
		return
	}
	defer repo.Delete(data.Key)

	result, err := repo.Get(data.Key)
	if err != nil {
		t.Errorf("failed while deleting data: %v", err)
		return
	}

	if result.Key != data.Key || result.Value != data.Value {
		t.Errorf("wrong data. expected %v got %v", data, result)
		return
	}
}
