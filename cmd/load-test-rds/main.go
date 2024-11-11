package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jvdbc/load-test-rds/internal/adapters"
	"github.com/jvdbc/load-test-rds/internal/models"
	"github.com/jvdbc/load-test-rds/internal/repositories"
	"github.com/jvdbc/load-test-rds/internal/services"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hostname",
				Value: "localhost",
				Usage: "message before counter",
			},
			&cli.StringFlag{
				Name:  "database",
				Value: "postgres",
				Usage: "name of the database",
			},
			&cli.StringFlag{
				Name:  "username",
				Value: "admin",
				Usage: "name of user",
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "poney",
				Usage: "password (use env PG_PASSWORD)",
			},
			&cli.UintFlag{
				Name:  "port",
				Value: 5432,
				Usage: "instance port",
			},
			// &cli.DurationFlag{
			// 	Name:  "frequency",
			// 	Value: time.Duration(1 * time.Second),
			// 	Usage: "elasped time between message",
			// },
		},
		Name:  "load-test-rds",
		Usage: "Connect to postgres",
		Action: func(ctx *cli.Context) error {
			hostname := ctx.String("hostname")
			port := ctx.Uint("port")
			database := ctx.String("database")
			username := ctx.String("username")
			password := ctx.String("password")

			connStr := models.NewConnectionString(hostname, database, port, username, password)
			db, err := open(connStr.String())
			if err != nil {
				return err
			}
			defer db.Close(context.Background())

			agentId := uint(1)

			orderAdapter := adapters.NewPostgresAdapter[models.Order](db)
			orderRepo := repositories.NewOrdersRepository(orderAdapter)
			orderService := services.NewOrdersService(orderRepo)

			return orderService.InsertNewOrderAndPrintAll(agentId)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func open(connection string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return db, nil
}

// func open2(connection string) (*sql.DB, error) {
// 	db, err := sql.Open("pgx", connection)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to connect to database: %w", err)
// 	}

// 	return db, nil
// }
