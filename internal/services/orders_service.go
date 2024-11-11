package services

import (
	"fmt"
	"os"

	"github.com/jvdbc/load-test-rds/internal/repositories"
)

type OrdersService struct {
	repo *repositories.OrdersRepository
}

func NewOrdersService(ordersRepo *repositories.OrdersRepository) *OrdersService {
	return &OrdersService{repo: ordersRepo}
}

func (me OrdersService) InsertNewOrderAndPrintAll(agentId uint) error {
	ordersCount, err := me.repo.Count(agentId)
	if err != nil {
		return err
	}

	_, err = me.repo.Insert(fmt.Sprintf("order %d from agent %d", ordersCount+1, agentId), agentId)
	if err != nil {
		return err
	}

	orders, err := me.repo.List(agentId)
	if err != nil {
		return err
	}

	for _, o := range orders {
		fmt.Fprintf(os.Stdout, "%s\n", o.String())
	}

	return nil
}
