package app

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolManager struct {
	Pools map[int64]*pgxpool.Pool
	mu    sync.RWMutex
}

func NewPoolManager() *PoolManager {
	return &PoolManager{
		Pools: make(map[int64]*pgxpool.Pool),
	}
}

func (pm *PoolManager) AddPool(id int64, connString string) (*pgxpool.Pool, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	pm.Pools[id] = pool
	return pool, nil
}

func (pm *PoolManager) GetPool(id int64) (*pgxpool.Pool, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	pool, exists := pm.Pools[id]
	return pool, exists
}

func (pm *PoolManager) DeletePool(id int64) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pool, exists := pm.Pools[id]
	if !exists {
		return errors.New("connection does not exist")
	}

	// Close the pool before deleting it
	pool.Close()

	// Remove the pool from the map
	delete(pm.Pools, id)

	return nil
}
