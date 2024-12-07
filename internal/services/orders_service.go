package services

import (
	"fmt"
	"os"
	"time"

	"github.com/jvdbc/load-test-rds/internal/repositories"
	"github.com/jvdbc/load-test-rds/tools"
)

type OrderWorker struct {
	id     uint
	ticker *time.Ticker
	repo   *repositories.OrdersRepository
}

func NewOrderWorker(agentId uint, every time.Duration, orderRepository *repositories.OrdersRepository) *OrderWorker {
	return &OrderWorker{id: agentId, ticker: time.NewTicker(every), repo: orderRepository}
}

func (oa OrderWorker) Id() uint {
	return oa.id
}

func (oa OrderWorker) StartInsert(begin uint) error {
	count := begin
	for range oa.ticker.C {
		count++
		_, err := oa.repo.Insert(fmt.Sprintf("order %d from agent %d", count, oa.id), oa.id)
		if err != nil {
			return fmt.Errorf("error in StartInsert: %w for iteration: %d", err, count)
		}
	}
	return nil
}

func (oa OrderWorker) Stop() {
	if oa.ticker != nil {
		oa.ticker.Stop()
	}
}

func (oa OrderWorker) StartPrintAll() error {
	for range oa.ticker.C {
		orders, err := oa.repo.List(oa.id)
		if err != nil {
			return err
		}

		tools.ExecClear(os.Stdout)

		for _, o := range orders {
			fmt.Fprintf(os.Stdout, "%s\n", o.String())
		}
	}
	return nil
}
