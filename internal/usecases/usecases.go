package usecases

import (
	"errors"
	"fmt"

	"github.com/vvjke314/vk-test-03-2025/internal/custom_errors"
	"github.com/vvjke314/vk-test-03-2025/internal/entities"
)

// repository defines the interface for data access operations
type repository interface {
	Insert(entities.VaultItem) error
	Update(key string, value string) error
	Delete(key string) error
	Get(key string) (entities.VaultItem, error)
	KeyExists(key string) (bool, error)
}

// KeyValueUseCase implements business logic for key-value operations
type KeyValueUseCase struct {
	repo repository
}

// NewKeyValueUseCase creates a new instance of KeyValueUseCase
func NewKeyValueUseCase(r repository) *KeyValueUseCase {
	return &KeyValueUseCase{
		repo: r,
	}
}

// InsertValue adds a new key-value pair after validation
func (uc *KeyValueUseCase) InsertValue(item entities.VaultItem) error {
	if item.Key == "" {
		return errors.New("key cannot be empty")
	}
	if item.Value == "" {
		return errors.New("value cannot be empty")
	}

	// Check if key exists
	exists, err := uc.repo.KeyExists(item.Key)
	if err != nil {
		return custom_errors.ErrKeyNotExists
	}
	if exists {
		return fmt.Errorf("key '%s' already exists", item.Key)
	}

	// Insert new record
	if err := uc.repo.Insert(item); err != nil {
		return fmt.Errorf("failed to insert value: %w", err)
	}

	return nil
}

// UpdateValue modifies an existing key-value pair
func (uc *KeyValueUseCase) UpdateValue(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if value == "" {
		return errors.New("value cannot be empty")
	}

	// Check if key exists
	exists, err := uc.repo.KeyExists(key)
	if err != nil {
		return fmt.Errorf("failed to check key existence: %w", err)
	}
	if !exists {
		return custom_errors.NewKeyNotExistsError(key)
	}

	// Update the record
	if err := uc.repo.Update(key, value); err != nil {
		return fmt.Errorf("failed to update value: %w", err)
	}

	return nil
}

// DeleteRow removes a key-value pair by key
func (uc *KeyValueUseCase) DeleteRow(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	// Check if key exists
	exists, err := uc.repo.KeyExists(key)
	if err != nil {
		return fmt.Errorf("failed to check key existence: %w", err)
	}
	if !exists {
		return custom_errors.NewKeyNotExistsError(key)
	}

	// Delete the record
	if err := uc.repo.Delete(key); err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}

	return nil
}

// Get retrieves a value by key
func (uc *KeyValueUseCase) Get(key string) (entities.VaultItem, error) {
	if key == "" {
		return entities.VaultItem{}, errors.New("key cannot be empty")
	}

	// Check if key exists
	exists, err := uc.repo.KeyExists(key)
	if err != nil {
		return entities.VaultItem{}, fmt.Errorf("failed to check key existence: %w", err)
	}
	if !exists {
		return entities.VaultItem{}, custom_errors.NewKeyNotExistsError(key)
	}

	item, err := uc.repo.Get(key)
	if err != nil {
		return entities.VaultItem{}, fmt.Errorf("failed to get value: %w", err)
	}

	return item, nil
}
