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
)

type TnRepository struct {
	dialer tarantool.NetDialer
	opts   tarantool.Opts
	conn   *tarantool.Connection
	config *config.TnRepoConfig
}

func NewTnRepository() *TnRepository {
	return &TnRepository{}
}

func (trepo *TnRepository) Init(ctx context.Context, cfg *config.TnRepoConfig) error {
	trepo.config = cfg
	trepo.dialer = tarantool.NetDialer{
		Address:  fmt.Sprintf("127.0.0.1:%s", trepo.config.Port),
		User:     cfg.Username,
		Password: cfg.Pass,
	}

	fmt.Println(trepo.dialer)

	var timeout time.Duration = time.Second * 5
	// if deadline, ok := ctx.Deadline(); !ok {
	// 	timeout = time.Second * 5 // if no timeout set
	// } else {
	// 	timeout = time.Until(deadline) // sync timeout with context
	// }

	trepo.opts = tarantool.Opts{
		Timeout: timeout,
	}

	var err error
	trepo.conn, err = tarantool.Connect(ctx, trepo.dialer, trepo.opts)
	if err != nil {
		return fmt.Errorf("error occured while connecting to the tarantool instance: %w", err)
	}

	return nil
}

func (trepo *TnRepository) Close() {
	trepo.conn.Close()
}

func (trepo *TnRepository) InsertData(i entities.VaultItem) error {
	tuple := []interface{}{i.Key, i.Value}
	future := trepo.conn.Do(tarantool.NewInsertRequest("vault").Tuple(tuple))
	rawResult, err := future.Get()
	if err != nil {
		return err
	}
	fmt.Println(rawResult...)
	return nil
}

func (trepo *TnRepository) GetAllData() ([]entities.VaultItem, error) {
	future := trepo.conn.Do(tarantool.NewSelectRequest("vault").Index("primary").Iterator(tarantool.IterAll).Limit(10000))

	resp, err := future.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to select data: %w", err)
	}

	var results []entities.VaultItem
	for _, tuple := range resp {
		tupleSlice, ok := tuple.([]interface{})
		if !ok || len(tupleSlice) < 2 {
			return nil, fmt.Errorf("invalid tuple format")
		}

		item := entities.VaultItem{
			Key:   fmt.Sprintf("%v", tupleSlice[0]), // Преобразуем interface{} в string
			Value: fmt.Sprintf("%v", tupleSlice[1]),
		}
		results = append(results, item)
	}

	return results, nil
}
