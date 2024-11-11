package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jvdbc/load-test-rds/internal/adapters"
	"github.com/jvdbc/load-test-rds/internal/models"
)

type OrdersRepository struct {
	db adapters.RelationalDatabaseTable[models.Order]
}

func NewOrdersRepository(adapter adapters.RelationalDatabaseTable[models.Order]) *OrdersRepository {
	return &OrdersRepository{db: adapter}
}

func (or OrdersRepository) List(agentId uint) ([]models.Order, error) {
	orders, err := or.db.Query(
		context.Background(),
		"SELECT id, content, created, agent_id FROM orders WHERE agent_id =$1 ORDER BY id",
		agentId)

	if err != nil {
		return nil, fmt.Errorf("list orders failed: %w", err)
	}

	return orders, nil
}

func (or OrdersRepository) Insert(content string, agentId uint) (*models.Order, error) {
	var newOrderId uint
	var newOrderCreated time.Time

	err := or.db.QueryRow(context.Background(), "INSERT INTO Orders (content, agent_id) VALUES($1, $2) RETURNING id, created", []any{content, agentId}, &newOrderId, &newOrderCreated)
	if err != nil {
		return nil, fmt.Errorf("insert into orders failed: %w", err)
	}

	return models.NewOrder(newOrderId, content, newOrderCreated, agentId), nil
}

func (or OrdersRepository) Count(agentId uint) (uint, error) {
	var count uint

	err := or.db.QueryRow(context.Background(), "SELECT COUNT(id) FROM orders WHERE agent_id =$1", []any{agentId}, &count)
	if err != nil {
		return count, fmt.Errorf("count orders failed: %w", err)
	}

	return count, nil
}
