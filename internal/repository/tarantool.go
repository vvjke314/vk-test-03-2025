package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
	"github.com/vvjke314/vk-test-03-2025/config"
	"github.com/vvjke314/vk-test-03-2025/internal/entities"
	"github.com/vvjke314/vk-test-03-2025/internal/logger"
)

// TnRepository represents Tarantool database repository
type TnRepository struct {
	dialer tarantool.NetDialer   // Connection dialer configuration
	opts   tarantool.Opts        // Connection options
	conn   *tarantool.Connection // Active connection
	config *config.TnRepoConfig  // Repository configuration
	logger logger.Logger         // Logger instance
}

// NewTnRepository creates new Tarantool repository instance
func NewTnRepository() *TnRepository {
	return &TnRepository{}
}

// Init initializes Tarantool connection with provided configuration
func (trepo *TnRepository) Init(ctx context.Context, cfg *config.TnRepoConfig, l logger.Logger) error {
	trepo.logger = l
	trepo.logger.Info("establishing connection with tarantool")
	trepo.config = cfg
	trepo.dialer = tarantool.NetDialer{
		Address:  fmt.Sprintf("tarantool:%s", trepo.config.Port),
		User:     cfg.Username,
		Password: cfg.Pass,
	}

	fmt.Println(trepo.dialer)

	// Set connection timeout
	var timeout time.Duration = time.Second * 5
	if deadline, ok := ctx.Deadline(); !ok {
		timeout = time.Second * 5 // Default timeout if not set in context
	} else {
		timeout = time.Until(deadline) // Use context deadline
	}

	trepo.opts = tarantool.Opts{
		Timeout: timeout,
	}

	// Establish connection
	var err error
	trepo.conn, err = tarantool.Connect(ctx, trepo.dialer, trepo.opts)
	if err != nil {
		err = fmt.Errorf("error occured while connecting to the tarantool instance: %w", err)
		trepo.logger.Error(err.Error())
		return err
	}

	trepo.logger.Info("tarantool client successfuly initialized")
	return nil
}

// Close terminates Tarantool connection
func (trepo *TnRepository) Close() {
	trepo.logger.Info("connection successfully closed")
	trepo.conn.Close()
}

// InsertData inserts new key-value pair into vault space
func (trepo *TnRepository) Insert(i entities.VaultItem) error {
	tuple := []interface{}{i.Key, i.Value}
	future := trepo.conn.Do(tarantool.NewInsertRequest("vault").Tuple(tuple))
	_, err := future.Get()
	if err != nil {
		trepo.logger.Error(err.Error())
		return err
	}
	trepo.logger.Info(fmt.Sprintf("inserted value with %s: %s values", i.Key, i.Value))
	return nil
}

// GetAllData retrieves all records from vault space
func (trepo *TnRepository) GetAllData() ([]entities.VaultItem, error) {
	trepo.logger.Info("scanning all data")
	future := trepo.conn.Do(tarantool.NewSelectRequest("vault").Index("primary").Iterator(tarantool.IterAll).Limit(10000))

	resp, err := future.Get()
	if err != nil {
		err = fmt.Errorf("failed to select data: %w", err)
		trepo.logger.Error(err.Error())
		return nil, err
	}

	// Convert response to VaultItem slice
	var results []entities.VaultItem
	for _, tuple := range resp {
		tupleSlice, ok := tuple.([]interface{})
		if !ok || len(tupleSlice) < 2 {
			err = fmt.Errorf("invalid tuple format")
			trepo.logger.Error(err.Error())
			return nil, err
		}

		item := entities.VaultItem{
			Key:   fmt.Sprintf("%v", tupleSlice[0]),
			Value: fmt.Sprintf("%v", tupleSlice[1]),
		}
		results = append(results, item)
	}

	trepo.logger.Info("scanned all data")
	return results, nil
}

// KeyExists checks if key exists in vault space
func (trepo *TnRepository) KeyExists(key string) (bool, error) {
	trepo.logger.Info(fmt.Sprintf("checking for key (%s) existence", key))
	var resp []entities.VaultItem
	err := trepo.conn.Do(tarantool.NewCallRequest(
		"key_check").Args([]interface{}{key})).GetTyped(&resp)

	if err != nil {
		err = fmt.Errorf("count failed: %w", err)
		trepo.logger.Error(err.Error())
		return false, err
	}
	trepo.logger.Info(fmt.Sprintf("key %s existence is %v", key, len(resp) > 0))
	return len(resp) > 0, nil
}

// Delete removes record with specified key from vault space
func (trepo *TnRepository) Delete(key string) error {
	trepo.logger.Info(fmt.Sprintf("deleting row with %s key", key))
	resp, err := trepo.conn.Do(tarantool.NewDeleteRequest("vault").Key([]interface{}{key})).Get()
	if err != nil {
		err = fmt.Errorf("delete failed: %w", err)
		trepo.logger.Error(err.Error())
		return err
	}

	trepo.logger.Info(fmt.Sprintf("successfully deleted %v", resp))
	return nil
}

// Update modifies value for existing key in vault space
func (trepo *TnRepository) Update(key string, value string) error {
	resp, err := trepo.conn.Do(
		tarantool.NewUpdateRequest("vault").
			Key([]interface{}{key}).
			Operations(tarantool.NewOperations().Assign(1, value))).Get()
	if err != nil {
		err = fmt.Errorf("update failed: %w", err)
		trepo.logger.Error(err.Error())
		return err
	}

	trepo.logger.Info(fmt.Sprintf("successfully updated %v", resp))
	return nil
}

// Get retrieves single record by key from vault space
func (trepo *TnRepository) Get(key string) (entities.VaultItem, error) {
	trepo.logger.Info(fmt.Sprintf("searching for row with %s key", key))
	var resp []entities.VaultItem
	err := trepo.conn.Do(tarantool.NewCallRequest(
		"key_check").Args([]interface{}{key})).GetTyped(&resp)

	if err != nil {
		err = fmt.Errorf("count failed: %w", err)
		trepo.logger.Error(err.Error())
		return entities.VaultItem{}, err
	}

	trepo.logger.Info(fmt.Sprintf("successfully got row %v with %s key", resp[0], key))
	return resp[0], nil
}
